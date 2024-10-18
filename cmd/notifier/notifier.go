package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/robertobadjio/tgtime-notifier/internal/aggregator"
	"github.com/robertobadjio/tgtime-notifier/internal/api_pb"
	"github.com/robertobadjio/tgtime-notifier/internal/background"
	"github.com/robertobadjio/tgtime-notifier/internal/config"
	kafkaLib "github.com/robertobadjio/tgtime-notifier/internal/kafka"
	"github.com/robertobadjio/tgtime-notifier/internal/notifier/telegram"
	"github.com/robertobadjio/tgtime-notifier/internal/notifier/telegram/command"
)

const checkSecondsInOffice = 10

func main() {
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	_ = logger.Log("msg", "Start service")

	ctx := context.Background()

	tgConfig, err := config.NewTelegramBotConfig()
	if err != nil {
		_ = logger.Log("config", "init", "type", "telegram", "error", err)
		os.Exit(1)
	}

	tgNotifier, err := telegram.NewTelegramNotifier(logger, tgConfig)
	if err != nil {
		_ = logger.Log("telegram", "init", "type", "bot", "error", err)
		os.Exit(1)
	}

	kafkaConfig, err := config.NewKafkaConfig()
	if err != nil {
		_ = logger.Log("kafka", "init", "type", "conn", "error", err)
		os.Exit(1)
	}

	startCheckInOffice(ctx, logger, tgNotifier, kafkaConfig.GetAddresses())
	startCheckPreviousDayInfo(ctx, logger, tgNotifier)

	updates := tgNotifier.GetBot().ListenForWebhook("/" + tgConfig.GetWebhookPath())
	go func() {
		srv := &http.Server{
			//Addr:         ":8441", // TODO: const
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
	logger log.Logger,
	tgNotifier *telegram.Notifier,
	addresses []string,
) {
	f := func() {
		kafka := kafkaLib.NewKafka(logger, addresses)
		err := kafka.ConsumeInOffice(ctx, tgNotifier)
		if err != nil {
			_ = logger.Log("kafka", "consume", "type", "in office message", "msg", err)
		}
	}
	bc := background.NewBackground(time.Duration(checkSecondsInOffice)*time.Second, f)
	bc.Start()
}

// TODO: Refactoring
func startCheckPreviousDayInfo(
	ctx context.Context,
	logger log.Logger,
	tn *telegram.Notifier,
) {
	f := func() {
		tgTimeAggregatorConfig, err := config.NewTgTimeAggregatorConfig()
		if err != nil {
			//return fmt.Errorf("error loading config: %w", err)
			_ = logger.Log("kafka", "consume", "type", "in office message", "msg", err)
			return
		}

		tgTimeAPIConfig, err := config.NewTgTimeAPIConfig()
		if err != nil {
			//return fmt.Errorf("error loading config: %w", err)
			_ = logger.Log("kafka", "consume", "type", "in office message", "msg", err)
			return
		}

		aggregatorClient := aggregator.NewClient(tgTimeAggregatorConfig, logger)

		timeSummaryResponse, err := aggregatorClient.GetTimeSummary(
			ctx,
			"",
			getPreviousDate("Europe/Moscow").Format("2006-01-02"),
		)
		if err != nil {
			//return fmt.Errorf("error getting time summary: %w", err)
			_ = logger.Log("kafka", "consume", "type", "in office message", "msg", err)
			return
		}
		if len(timeSummaryResponse.Summary) == 0 {
			//return nil
			_ = logger.Log("kafka", "consume", "type", "in office message", "msg", err)
			return
		}

		apiClient := api_pb.NewClient(tgTimeAPIConfig, logger)

		for _, summaryByUser := range timeSummaryResponse.Summary {
			user, err := apiClient.GetUserByMacAddress(ctx, summaryByUser.MacAddress)
			if err != nil {
				_ = logger.Log("kafka", "consume", "type", "in office message", "msg", err)
				continue
				//return fmt.Errorf("error getting user by telegram id: %w", err)
			}

			hours, minutes := command.SecondsToHM(summaryByUser.GetSeconds())

			breaksString := command.BreaksToString(command.BuildBreaks(summaryByUser.Breaks))
			message := fmt.Sprintf(
				"Вчера Вы были в офисе с %s до %s\nУчтенное время %d ч. %d м.",
				command.SecondsToTime(summaryByUser.GetSecondsStart()).Format("15:04"),
				command.SecondsToTime(summaryByUser.GetSecondsEnd()).Format("15:04"),
				hours,
				minutes,
			)

			if "" != breaksString {
				message += fmt.Sprintf("\nПерерывы %s\n", breaksString)
			}

			_, err = tn.GetBot().Send(tn.SetKeyboard(tgbotapi.NewMessage(user.GetUser().TelegramId, message)))
			if err != nil {
				_ = logger.Log("kafka", "consume", "type", "in office message", "msg", err)
			}
		}
	}
	// TODO: В 12 дня посылать информацию о предыдущем дне
	bc2 := background.NewBackground(time.Duration(60)*time.Second, f)
	bc2.Start()
}

func getPreviousDate(location string) time.Time {
	moscowLocation, _ := time.LoadLocation(location)
	return time.Now().AddDate(0, 0, -1).In(moscowLocation)
}
