package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	_ "github.com/lib/pq"
	"github.com/robertobadjio/tgtime-notifier/internal/background"
	"github.com/robertobadjio/tgtime-notifier/internal/config"
	kafkaLib "github.com/robertobadjio/tgtime-notifier/internal/kafka"
	"github.com/robertobadjio/tgtime-notifier/internal/notifier/telegram"
)

const checkSecondsInOffice = 10

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
		srv := &http.Server{
			Addr:         ":8441", // TODO: const
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
		err := srv.ListenAndServe()
		if err != nil {
			_ = logger.Log("telegram", "updates", "type", "serve", "msg", err)
		}
	}()

	for update := range updates {
		if update.Message == nil {
			continue
		}
		err := tgNotifier.SendMessageCommand(ctx, update)
		if err != nil {
			_ = logger.Log("telegram", "updates", "type", "info", "err", err)
		}
	}
}

func startCheckInOffice(
	ctx context.Context,
	cfg *config.Config,
	logger log.Logger,
	tgNotifier *telegram.Notifier,
) {
	f := func() {
		kafka := kafkaLib.NewKafka(logger, cfg.KafkaHost, cfg.KafkaPort)
		err := kafka.ConsumeInOffice(ctx, tgNotifier)
		if err != nil {
			_ = logger.Log("kafka", "consume", "type", "in office message", "msg", err)
		}
	}
	bc := background.NewBackground(time.Duration(checkSecondsInOffice)*time.Second, f)
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
