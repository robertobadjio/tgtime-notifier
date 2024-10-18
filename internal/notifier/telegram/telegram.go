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
func NewTelegramNotifier(
	logger log.Logger,
	tgConfig config.TelegramBotConfig,
) (*Notifier, error) {
	bot, err := initTelegramBot(tgConfig.GetToken())
	if err != nil {
		return nil, fmt.Errorf("error creating telegram bot: %w", err)
	}
	_ = logger.Log("notifier", "telegram", "name", bot.Self.UserName, "msg", "authorized on account")

	err = setWebhook(bot, tgConfig.GetWebhookLink())
	if err != nil {
		return nil, fmt.Errorf("error setting telegram webhook: %w", err)
	}
	_ = logger.Log(
		"notifier", "telegram",
		"name", bot.Self.UserName,
		"msg", "setting webhook",
		"url", tgConfig.GetWebhookLink(),
	)

	return &Notifier{logger: logger, Bot: bot}, nil
}

// GetBot Получение Telegram Bot API
func (tn *Notifier) GetBot() *tgbotapi.BotAPI {
	return tn.Bot
}

func initTelegramBot(token string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("init telegram bot: %w", err)
	}

	return bot, nil
}

func setWebhook(bot *tgbotapi.BotAPI, webhookPath string) error {
	_, err := bot.SetWebhook(tgbotapi.NewWebhook(webhookPath))
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

// SetKeyboard ???
func (*Notifier) SetKeyboard(message tgbotapi.MessageConfig) tgbotapi.MessageConfig {
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
	_, err = tn.Bot.Send(tn.SetKeyboard(tgbotapi.NewMessage(
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
	_, err = tn.Bot.Send(tn.SetKeyboard(tgbotapi.NewMessage(
		telegramID,
		stringMessage,
	)))
	if err != nil {
		return fmt.Errorf("error send welcome message: %w", err)
	}

	return nil
}

func (tn *Notifier) SendPreviousDayInfoMessage(ctx context.Context) error {
	// TODO: ?!

	return nil
}
