package main

import (
	"context"
	"github.com/go-kit/kit/log"
	_ "github.com/lib/pq"
	"net/http"
	"os"
	"tgtime-notifier/internal/background"
	"tgtime-notifier/internal/config"
	kafkaLib "tgtime-notifier/internal/kafka"
	"tgtime-notifier/internal/notifier/telegram"
	"time"
)

func main() {
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))

	_ = logger.Log("msg", "Start service")

	cfg := config.New()
	ctx := context.Background()
	tgNotifier := telegram.NewTelegramNotifier(logger)

	f := func() {
		kafka := kafkaLib.NewKafka(logger, cfg.KafkaHost, cfg.KafkaPort)
		err := kafka.ConsumeInOffice(ctx, tgNotifier)
		if err != nil {
			_ = logger.Log("kafka", "consume", "type", "in office message", "msg", err)
		}
	}
	bc := background.NewBackground(time.Duration(10)*time.Second, f)
	bc.Start()

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

	updates := tgNotifier.GetBot().ListenForWebhook("/" + cfg.WebHookPath)
	go http.ListenAndServe(":8441", nil) // TODO: Handle error
	for update := range updates {
		tgNotifier.Info(ctx, update)
	}
}
