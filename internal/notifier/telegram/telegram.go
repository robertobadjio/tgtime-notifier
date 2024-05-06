package telegram

import (
	"fmt"
	"github.com/go-kit/kit/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
	"tgtime-notifier/internal/config"
	"tgtime-notifier/internal/notifier"
	"time"
)

const (
	buttonWorkingTime              = "⏳ Рабочее время"
	buttonStatCurrentWorkingPeriod = "🗓 Статистика за рабочий период"
	buttonStart                    = "/start"
)

type TelegramNotifier struct {
	logger log.Logger
	bot    *tgbotapi.BotAPI
}

func NewTelegramNotifier(logger log.Logger) *TelegramNotifier {
	bot, err := initTelegramBot()
	if err != nil {
		panic(err)
	}

	_ = logger.Log("notifier", "telegram", "name", bot.Self.UserName, "msg", "authorized on account")

	err = setWebhook(bot)
	if err != nil {
		panic(err)
	}

	return &TelegramNotifier{logger: logger, bot: bot}
}

func (t *TelegramNotifier) GetBot() *tgbotapi.BotAPI {
	return t.bot
}

func initTelegramBot() (*tgbotapi.BotAPI, error) {
	cfg := config.New()
	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, fmt.Errorf("init telegram bot: %w", err)
	}

	return bot, nil
}

func setWebhook(bot *tgbotapi.BotAPI) error {
	cfg := config.New()
	_, err := bot.SetWebhook(tgbotapi.NewWebhook(cfg.WebHookLink))
	if err != nil {
		return fmt.Errorf("set telegram webhook: %w", err)
	}

	//info, err := bot.GetWebhookInfo()
	_, err = bot.GetWebhookInfo()
	if err != nil {
		return fmt.Errorf("get telegram webhook info: %w", err)
	}
	/*if info.LastErrorDate != 0 {
		return fmt.Errorf("telegram callback failed: %s", info.LastErrorMessage)
	}*/

	return nil
}

// Keyboard
func (t *TelegramNotifier) setKeyboard(message tgbotapi.MessageConfig) tgbotapi.MessageConfig {
	message.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(buttonWorkingTime),
			tgbotapi.NewKeyboardButton(buttonStatCurrentWorkingPeriod),
		),
	)
	return message
}

func buildBreaks(breaks []*notifier.Break) []string {
	var output []string
	for _, item := range breaks {
		beginTime := time.Unix(item.StartTime, 0)
		endTime := time.Unix(item.EndTime, 0)
		output = append(
			output,
			fmt.Sprintf("%s - %s", beginTime.Format("15:04"), endTime.Format("15:04")))
	}

	return output
}

func breaksToString(breaks []string) string {
	return strings.Join(breaks, ", ")
}
