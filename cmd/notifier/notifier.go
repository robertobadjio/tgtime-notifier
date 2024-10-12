package main

import (
	"context"
	"github.com/go-kit/kit/log"
	_ "github.com/lib/pq"
	"net/http"
	"os"
	"tgtime-notifier/internal/aggregator"
	"tgtime-notifier/internal/api_pb"
	"tgtime-notifier/internal/background"
	"tgtime-notifier/internal/config"
	kafkaLib "tgtime-notifier/internal/kafka"
	"tgtime-notifier/internal/notifier/telegram"
	"time"
)

const CheckSecondsInOffice = 10

func main() {
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	_ = logger.Log("msg", "Start service")

	cfg := config.New()
	ctx := context.Background()
	tgNotifier := telegram.NewTelegramNotifier(logger)

	startCheckInOffice(ctx, cfg, logger, tgNotifier)
	startCheckPreviousDayInfo()

	updates := tgNotifier.GetBot().ListenForWebhook("/" + cfg.WebHookPath)
	go func() {
		err := http.ListenAndServe(":1124", nil) // TODO: const
		if err != nil {
			_ = logger.Log("telegram", "updates", "type", "serve", "msg", err)
		}
	}()

	aggregatorClient := aggregator.NewClient(*cfg, logger)
	apiClient := api_pb.NewClient(*cfg, logger)
	for update := range updates {
		err := tgNotifier.Info(ctx, update, aggregatorClient, apiClient)
		if err != nil {
			_ = logger.Log("telegram", "updates", "type", "info", "err", err)
		}
	}
}

func startCheckInOffice(
	ctx context.Context,
	cfg *config.Config,
	logger log.Logger,
	tgNotifier *telegram.TelegramNotifier,
) {
	f := func() {
		kafka := kafkaLib.NewKafka(logger, cfg.KafkaHost, cfg.KafkaPort)
		err := kafka.ConsumeInOffice(ctx, tgNotifier)
		if err != nil {
			_ = logger.Log("kafka", "consume", "type", "in office message", "msg", err)
		}
	}
	bc := background.NewBackground(time.Duration(CheckSecondsInOffice)*time.Second, f)
	bc.Start()
}

func startCheckPreviousDayInfo() {
	/*f2 := func() {
		kafka := kafkaLib.NewKafka(logger, cfg.KafkaHost, cfg.KafkaPort)
		err := kafka.ConsumePreviousDayInfo(ctx, tgNotifier)
		if err != nil {
			_ = logger.Log("msg", "consume previous day info failed", "err", err)
		}
	}
	// TODO: В 12 дня посылать информацию о предыдущем дне
	bc2 := background.NewBackground(time.Duration(60)*time.Second, f2)
	bc2.Start()*/
}
