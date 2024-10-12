package telegram

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"tgtime-notifier/internal/config"
)

type TelegramNotifier struct {
	logger log.Logger
	Bot    *tgbotapi.BotAPI
}

func NewTelegramNotifier(logger log.Logger) *TelegramNotifier {
	bot, err := initTelegramBot()
	if err != nil {
		panic(err)
	}

	_ = logger.Log("notifier", "telegram", "name", bot.Self.UserName, "msg", "authorized on account")
	_ = logger.Log("notifier", "telegram", "name", bot.Self.UserName, "msg", "setting webhook...")
	err = setWebhook(bot)
	if err != nil {
		panic(err)
	}

	return &TelegramNotifier{logger: logger, Bot: bot}
}

func (tn *TelegramNotifier) GetBot() *tgbotapi.BotAPI {
	return tn.Bot
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
func (*TelegramNotifier) SetKeyboard(message tgbotapi.MessageConfig) tgbotapi.MessageConfig {
	message.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(string(buttonWorkingTime)),
			tgbotapi.NewKeyboardButton(string(buttonStatCurrentWorkingPeriod)),
		),
	)
	return message
}

func (tn *TelegramNotifier) SendMessageCommand(ctx context.Context, update tgbotapi.Update) error {
	command := NewCommand(MessageType(update.Message.Text), int64(update.Message.From.ID))
	messageHandler := TypeMessage{Message: command}
	stringMessage, err := messageHandler.Handle(ctx)
	_, err = tn.Bot.Send(tn.SetKeyboard(tgbotapi.NewMessage(
		int64(update.Message.From.ID),
		stringMessage,
	)))

	return fmt.Errorf("error send telegram message: %w", err)
}

func (tn *TelegramNotifier) SendWelcomeMessage(_ context.Context, telegramId int64) error {
	_, err := tn.Bot.Send(tn.SetKeyboard(tgbotapi.NewMessage(
		telegramId,
		"Вы пришли в офис",
	)))

	return fmt.Errorf("error send telegram welcome message: %w", err)
}
