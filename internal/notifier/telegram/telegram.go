package telegram

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/robertobadjio/tgtime-notifier/internal/config"
)

// Notifier Telegram-нотификатор
type Notifier struct {
	logger log.Logger
	Bot    *tgbotapi.BotAPI
}

// NewTelegramNotifier Конструктор для создания Telegram-нотификатора
func NewTelegramNotifier(logger log.Logger) *Notifier {
	bot, err := initTelegramBot()
	if err != nil {
		panic(err)
	}

	_ = logger.Log("notifier", "telegram", "name", bot.Self.UserName, "msg", "authorized on account")
	/*_ = logger.Log("notifier", "telegram", "name", bot.Self.UserName, "msg", "setting webhook...")
	err = setWebhook(bot)
	if err != nil {
		panic(err)
	}*/

	return &Notifier{logger: logger, Bot: bot}
}

// GetBot Получение Telegram Bot API
func (tn *Notifier) GetBot() *tgbotapi.BotAPI {
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

func (*Notifier) setKeyboard(message tgbotapi.MessageConfig) tgbotapi.MessageConfig {
	message.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(string(buttonWorkingTime)),
			tgbotapi.NewKeyboardButton(string(buttonStatCurrentWorkingPeriod)),
		),
	)
	return message
}

// SendMessageCommand Метод для отправки сообщения в ответ на команду пользователя
func (tn *Notifier) SendMessageCommand(ctx context.Context, update tgbotapi.Update) error {
	command := NewCommand(MessageType(update.Message.Text), int64(update.Message.From.ID))
	stringMessage, err := command.GetMessage(ctx)
	if err != nil {
		return fmt.Errorf("error getting text message: %w", err)
	}
	_, err = tn.Bot.Send(tn.setKeyboard(tgbotapi.NewMessage(
		int64(update.Message.From.ID),
		stringMessage,
	)))

	if err != nil {
		return fmt.Errorf("error send telegram message: %w", err)
	}

	return nil
}

// SendWelcomeMessage Метод отправки приветственного сообщения по приходу в офис / на работу.
func (tn *Notifier) SendWelcomeMessage(ctx context.Context, telegramID int64) error {
	command := NewCommand(welcome, telegramID)
	stringMessage, err := command.GetMessage(ctx)
	if err != nil {
		return fmt.Errorf("error send welcome message: %w", err)
	}
	_, err = tn.Bot.Send(tn.setKeyboard(tgbotapi.NewMessage(
		telegramID,
		stringMessage,
	)))
	if err != nil {
		return fmt.Errorf("error send welcome message: %w", err)
	}

	return nil
}
