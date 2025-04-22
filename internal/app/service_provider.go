package app

import (
	"context"
	"net/http"

	TGBotAPI "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/robertobadjio/tgtime-notifier/internal/api"
	"github.com/robertobadjio/tgtime-notifier/internal/api/service/endpoints"
	"github.com/robertobadjio/tgtime-notifier/internal/api/service/transport"
	"github.com/robertobadjio/tgtime-notifier/internal/config"
	"github.com/robertobadjio/tgtime-notifier/internal/kafka"
	"github.com/robertobadjio/tgtime-notifier/internal/logger"
	"github.com/robertobadjio/tgtime-notifier/internal/metric"
	"github.com/robertobadjio/tgtime-notifier/internal/service/client/aggregator"
	"github.com/robertobadjio/tgtime-notifier/internal/service/client/api_pb"
	"github.com/robertobadjio/tgtime-notifier/internal/service/notifier/telegram"
	"github.com/robertobadjio/tgtime-notifier/internal/service/previous_day_info"
)

type serviceProvider struct {
	httpConfig config.HTTPConfig

	httpServiceHandler  http.Handler
	endpointsServiceSet endpoints.Set
	service             api.Service

	tgConfig   config.TelegramBotConfig
	tgBot      *TGBotAPI.BotAPI
	tgNotifier *telegram.TGNotifier

	kafkaConfig config.KafkaConfig
	kafka       *kafka.Kafka

	tgTimeAPIConfig config.TgTimeAPIConfig
	tgTimeAPIClient api_pb.Client

	tgTimeAggregatorConfig config.TgTimeAggregatorConfig
	tgTimeAggregatorClient aggregator.Client

	previousDayInfo previous_day_info.PreviousDayInfo

	promConfig      config.PromConfig
	pyroscopeConfig config.PyroscopeConfig
	metrics         *metric.Metrics

	os config.OS
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// HTTPServiceHandler ...
func (sp *serviceProvider) HTTPServiceHandler(ctx context.Context) http.Handler {
	if sp.httpServiceHandler == nil {
		sp.httpServiceHandler = transport.NewHTTPHandler(sp.EndpointsServiceSet(ctx))
	}

	return sp.httpServiceHandler
}

// EndpointsServiceSet ...
func (sp *serviceProvider) EndpointsServiceSet(ctx context.Context) endpoints.Set {
	sp.endpointsServiceSet = endpoints.NewEndpointSet(sp.APIService(ctx))

	return sp.endpointsServiceSet
}

// OS ...
func (sp *serviceProvider) OS() config.OS {
	if sp.os == nil {
		sp.os = config.NewOS()
	}

	return sp.os
}

// APIService ...
func (sp *serviceProvider) APIService(_ context.Context) api.Service {
	if sp.service == nil {
		sp.service = api.NewNotifierService()
	}

	return sp.service
}

// HTTPConfig ...
func (sp *serviceProvider) HTTPConfig() config.HTTPConfig {
	if sp.httpConfig == nil {
		httpConfig, err := config.NewHTTPConfig(sp.OS())
		if err != nil {
			logger.Fatal("di", "http", "error", err.Error())
		}

		sp.httpConfig = httpConfig
	}

	return sp.httpConfig
}

// PromConfig ...
func (sp *serviceProvider) PromConfig() config.PromConfig {
	if sp.promConfig == nil {
		promConfig, err := config.NewPromConfig(sp.OS())
		if err != nil {
			logger.Fatal("di", "prometheus", "error", err.Error())
		}

		sp.promConfig = promConfig
	}

	return sp.promConfig
}

// Metrics ...
func (sp *serviceProvider) Metrics() *metric.Metrics {
	if sp.metrics == nil {
		sp.metrics = metric.NewMetrics()
	}

	return sp.metrics
}

// PyroscopeConfig ...
func (sp *serviceProvider) PyroscopeConfig() config.PyroscopeConfig {
	if sp.pyroscopeConfig == nil {
		pyroscopeConfig, err := config.NewPyroscopeConfig(sp.OS())
		if err != nil {
			logger.Fatal("di", "pyroscope", "error", err.Error())
		}

		sp.pyroscopeConfig = pyroscopeConfig
	}

	return sp.pyroscopeConfig
}

// KafkaConfig ...
func (sp *serviceProvider) KafkaConfig() config.KafkaConfig {
	if sp.kafkaConfig == nil {
		kc, err := config.NewKafkaConfig(sp.OS())
		if err != nil {
			logger.Fatal("di", "http", "error", err.Error())
		}

		sp.kafkaConfig = kc
	}

	return sp.kafkaConfig
}

// TelegramConfig ...
func (sp *serviceProvider) TelegramConfig() config.TelegramBotConfig {
	if sp.tgConfig == nil {
		tgConfig, err := config.NewTelegramBotConfig(sp.OS())
		if err != nil {
			logger.Fatal("di", "tgConfig", "error", err.Error())
		}

		sp.tgConfig = tgConfig
	}

	return sp.tgConfig
}

// TGNotifier ...
func (sp *serviceProvider) TGNotifier() *telegram.TGNotifier {
	if sp.tgNotifier == nil {
		tgNC, err := telegram.NewTelegramNotifier(
			sp.TgBot(),
			sp.Metrics(),
			sp.TGTimeAPIClient(),
			sp.TGTimeAggregatorClient(),
		)
		if err != nil {
			logger.Fatal("di", "tgNotifier", "error", err.Error())
		}

		sp.tgNotifier = tgNC
	}

	return sp.tgNotifier
}

// TgBot ...
func (sp *serviceProvider) TgBot() *TGBotAPI.BotAPI {
	if sp.tgBot == nil {
		bot, err := TGBotAPI.NewBotAPI(sp.TelegramConfig().Token())
		if err != nil {
			logger.Fatal("di", "tgBot", "error", err.Error())
		}

		_, err = bot.SetWebhook(TGBotAPI.NewWebhook(sp.TelegramConfig().WebhookLink()))
		if err != nil {
			logger.Fatal("di", "tgBot", "type", "set telegram webhook", "error", err.Error())
		}

		//info, err := bot.GetWebhookInfo()
		_, err = bot.GetWebhookInfo()
		if err != nil {
			logger.Fatal("di", "tgBot", "type", "get telegram webhook", "error", err.Error())
		}
		/*if info.LastErrorDate != 0 {
			return fmt.Errorf("telegram callback failed: %s", info.LastErrorMessage)
		}*/

		sp.tgBot = bot
	}

	return sp.tgBot
}

// Kafka ...
func (sp *serviceProvider) Kafka() *kafka.Kafka {
	if sp.kafka == nil {
		sp.kafka = kafka.NewKafka(sp.KafkaConfig(), sp.TGNotifier(), sp.TGTimeAPIClient())
	}

	return sp.kafka
}

// TGTimeAPIConfig ...
func (sp *serviceProvider) TGTimeAPIConfig() config.TgTimeAPIConfig {
	if sp.tgTimeAPIConfig == nil {
		c, err := config.NewTgTimeAPIConfig(sp.OS())
		if err != nil {
			logger.Fatal("di", "tgTimeAPIConfig", "error", err.Error())
		}

		sp.tgTimeAPIConfig = c
	}

	return sp.tgTimeAPIConfig
}

// TGTimeAPIClient ...
func (sp *serviceProvider) TGTimeAPIClient() api_pb.Client {
	if sp.tgTimeAPIClient == nil {
		sp.tgTimeAPIClient = api_pb.NewClient(sp.TGTimeAPIConfig())
	}

	return sp.tgTimeAPIClient
}

// TGTimeAggregatorConfig ...
func (sp *serviceProvider) TGTimeAggregatorConfig() config.TgTimeAPIConfig {
	if sp.tgTimeAggregatorConfig == nil {
		c, err := config.NewTgTimeAggregatorConfig(sp.OS())
		if err != nil {
			logger.Fatal("di", "tgTimeAggregatorConfig", "error", err.Error())
		}

		sp.tgTimeAggregatorConfig = c
	}

	return sp.tgTimeAggregatorConfig
}

// TGTimeAggregatorClient ...
func (sp *serviceProvider) TGTimeAggregatorClient() aggregator.Client {
	if sp.tgTimeAggregatorClient == nil {
		sp.tgTimeAggregatorClient = aggregator.NewClient(sp.TGTimeAggregatorConfig())
	}

	return sp.tgTimeAggregatorClient
}

// PreviousDayInfo ...
func (sp *serviceProvider) PreviousDayInfo() previous_day_info.PreviousDayInfo {
	if sp.previousDayInfo == nil {
		sp.previousDayInfo = previous_day_info.NewPreviousDayInfo(
			sp.TGTimeAPIClient(),
			sp.TGTimeAggregatorClient(),
			sp.TGNotifier(),
		)
	}

	return sp.previousDayInfo
}
