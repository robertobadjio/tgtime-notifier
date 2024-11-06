package app

import (
	TGBotAPI "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/robertobadjio/tgtime-notifier/internal/aggregator"
	"github.com/robertobadjio/tgtime-notifier/internal/api_pb"
	"github.com/robertobadjio/tgtime-notifier/internal/config"
	"github.com/robertobadjio/tgtime-notifier/internal/kafka"
	"github.com/robertobadjio/tgtime-notifier/internal/logger"
	notifierI "github.com/robertobadjio/tgtime-notifier/internal/notifier"
	"github.com/robertobadjio/tgtime-notifier/internal/notifier/telegram"
	"github.com/robertobadjio/tgtime-notifier/internal/notifier/telegram/command"
	"github.com/robertobadjio/tgtime-notifier/internal/task"
	"github.com/robertobadjio/tgtime-notifier/internal/task/previous_day_info"
)

type serviceProvider struct {
	httpConfig config.HTTPConfig

	tgConfig         config.TelegramBotConfig
	tgBot            *TGBotAPI.BotAPI
	tgNotifier       notifierI.Notifier
	tgCommandFactory command.Factory

	kafkaConfig config.KafkaConfig
	kafka       *kafka.Kafka

	tgTimeAPIConfig config.TgTimeAPIConfig
	tgTimeAPIClient api_pb.Client

	tgTimeAggregatorConfig config.TgTimeAggregatorConfig
	tgTimeAggregatorClient aggregator.Client

	previousDayInfoTask task.Task
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (sp *serviceProvider) HTTPConfig() config.HTTPConfig {
	if sp.httpConfig == nil {
		httpConfig, err := config.NewHTTPConfig()
		if err != nil {
			logger.Fatal("di", "http", "error", err.Error())
		}

		sp.httpConfig = httpConfig
	}

	return sp.httpConfig
}

func (sp *serviceProvider) KafkaConfig() config.KafkaConfig {
	if sp.kafkaConfig == nil {
		kc, err := config.NewKafkaConfig()
		if err != nil {
			logger.Fatal("di", "http", "error", err.Error())
		}

		sp.kafkaConfig = kc
	}

	return sp.kafkaConfig
}

func (sp *serviceProvider) TelegramConfig() config.TelegramBotConfig {
	if sp.tgConfig == nil {
		tgConfig, err := config.NewTelegramBotConfig()
		if err != nil {
			logger.Fatal("di", "tgConfig", "error", err.Error())
		}

		sp.tgConfig = tgConfig
	}

	return sp.tgConfig
}

func (sp *serviceProvider) TGNotifier() notifierI.Notifier {
	if sp.tgNotifier == nil {
		tgNC, err := telegram.NewTelegramNotifier(sp.TgBot(), sp.TGCommandFactory())
		if err != nil {
			logger.Fatal("di", "tgNotifier", "error", err.Error())
		}

		sp.tgNotifier = tgNC
	}

	return sp.tgNotifier
}

func (sp *serviceProvider) TgBot() *TGBotAPI.BotAPI {
	if sp.tgBot == nil {
		bot, err := TGBotAPI.NewBotAPI(sp.tgConfig.GetToken())
		if err != nil {
			logger.Fatal("di", "tgBot", "error", err.Error())
		}

		sp.tgBot = bot
	}

	return sp.tgBot
}

func (sp *serviceProvider) Kafka() *kafka.Kafka {
	if sp.kafka == nil {
		sp.kafka = kafka.NewKafka(sp.KafkaConfig(), sp.TGNotifier(), sp.TGTimeAPIClient())
	}

	return sp.kafka
}

func (sp *serviceProvider) TGTimeAPIConfig() config.TgTimeAPIConfig {
	if sp.tgTimeAPIConfig == nil {
		c, err := config.NewTgTimeAPIConfig()
		if err != nil {
			logger.Fatal("di", "tgTimeAPIConfig", "error", err.Error())
		}

		sp.tgTimeAPIConfig = c
	}

	return sp.tgTimeAPIConfig
}

func (sp *serviceProvider) TGTimeAPIClient() api_pb.Client {
	if sp.tgTimeAPIClient == nil {
		sp.tgTimeAPIClient = api_pb.NewClient(sp.TGTimeAPIConfig())
	}

	return sp.tgTimeAPIClient
}

func (sp *serviceProvider) TGTimeAggregatorConfig() config.TgTimeAPIConfig {
	if sp.tgTimeAggregatorConfig == nil {
		c, err := config.NewTgTimeAggregatorConfig()
		if err != nil {
			logger.Fatal("di", "tgTimeAggregatorConfig", "error", err.Error())
		}

		sp.tgTimeAggregatorConfig = c
	}

	return sp.tgTimeAggregatorConfig
}

func (sp *serviceProvider) TGTimeAggregatorClient() aggregator.Client {
	if sp.tgTimeAggregatorClient == nil {
		sp.tgTimeAggregatorClient = aggregator.NewClient(sp.TGTimeAggregatorConfig())
	}

	return sp.tgTimeAggregatorClient
}

func (sp *serviceProvider) PreviousDayInfoTask() task.Task {
	if sp.previousDayInfoTask == nil {
		sp.previousDayInfoTask = previous_day_info.NewPreviousDayInfoTask(
			sp.TGTimeAPIClient(),
			sp.TGTimeAggregatorClient(),
			sp.TGNotifier(),
		)
	}

	return sp.previousDayInfoTask
}

func (sp *serviceProvider) TGCommandFactory() command.Factory {
	if sp.tgCommandFactory == nil {
		sp.tgCommandFactory = command.NewFactory(
			sp.TGTimeAPIClient(),
			sp.TGTimeAggregatorClient(),
		)
	}

	return sp.tgCommandFactory
}
