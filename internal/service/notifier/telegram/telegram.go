package telegram

import (
	"fmt"
	"time"

	TGBotAPI "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/robertobadjio/tgtime-notifier/internal/metric"
	notifierI "github.com/robertobadjio/tgtime-notifier/internal/service/notifier"
	command2 "github.com/robertobadjio/tgtime-notifier/internal/service/notifier/telegram/command"
)

// TgNotifier Telegram-нотификатор
type notifier struct {
	bot     *TGBotAPI.BotAPI
	factory command2.Factory
}

// ParamsUpdate ???
type ParamsUpdate struct {
	Update TGBotAPI.Update
}

// ParamsWorkingTime ???
type ParamsWorkingTime struct {
}

// ParamsPreviousDayInfo ???
type ParamsPreviousDayInfo struct {
	TelegramID   int64
	SecondsStart time.Time
	SecondsEnd   time.Time
	Hours        int64
	Minutes      int64
	Breaks       string
}

// ParamsWelcomeMessage ???
type ParamsWelcomeMessage struct {
	TelegramID int64
}

// NewTelegramNotifier Конструктор для создания Telegram-нотификатора
func NewTelegramNotifier(
	bot *TGBotAPI.BotAPI,
	factory command2.Factory,
) (notifierI.Notifier, error) {
	return &notifier{bot: bot, factory: factory}, nil
}

// Bot Получение Telegram Bot API
func (tn *notifier) Bot() *TGBotAPI.BotAPI {
	return tn.bot
}

func (tn *notifier) Factory() command2.Factory {
	return tn.factory
}

// SetKeyboard ???
func (tn *notifier) SetKeyboard(message TGBotAPI.MessageConfig) TGBotAPI.MessageConfig {
	message.ReplyMarkup = TGBotAPI.NewReplyKeyboard(
		TGBotAPI.NewKeyboardButtonRow(
			TGBotAPI.NewKeyboardButton(string(command2.ButtonWorkingTime)),
			TGBotAPI.NewKeyboardButton(string(command2.ButtonStatCurrentWorkingPeriod)),
		),
	)
	return message
}

// SendMessage ...
func (tn *notifier) SendMessage(text string, telegramID int64) error {
	_, err := tn.Bot().Send(tn.SetKeyboard(TGBotAPI.NewMessage(
		telegramID,
		text,
	)))
	if err != nil {
		return fmt.Errorf("error send message: %w", err)
	}

	metric.IncMessageCounter()

	return nil
}
