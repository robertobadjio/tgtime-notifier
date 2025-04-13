package telegram

import (
	"errors"
	"fmt"
	"time"

	TGBotAPI "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/robertobadjio/tgtime-notifier/internal/metric"
	notifierI "github.com/robertobadjio/tgtime-notifier/internal/service/notifier"
	command2 "github.com/robertobadjio/tgtime-notifier/internal/service/notifier/telegram/command"
)

// BotAPIInterface ...
type BotAPIInterface interface {
	Send(c TGBotAPI.Chattable) (TGBotAPI.Message, error)
}

// TgNotifier Telegram-нотификатор
type notifier struct {
	bot     BotAPIInterface
	factory command2.Factory
	metrics metric.Metrics
}

// ParamsUpdate ...
type ParamsUpdate struct {
	Update TGBotAPI.Update
}

// ParamsWorkingTime ...
type ParamsWorkingTime struct {
}

// ParamsPreviousDayInfo ...
type ParamsPreviousDayInfo struct {
	TelegramID   int64
	SecondsStart time.Time
	SecondsEnd   time.Time
	Hours        int64
	Minutes      int64
	Breaks       string
}

// ParamsWelcomeMessage ...
type ParamsWelcomeMessage struct {
	TelegramID int64
}

// NewTelegramNotifier Конструктор для создания Telegram-нотификатора
func NewTelegramNotifier(
	bot BotAPIInterface,
	factory command2.Factory,
	metrics metric.Metrics,
) (notifierI.Notifier, error) {
	if bot == nil {
		return nil, errors.New("telegram bot is nil")
	}

	if factory == nil {
		return nil, errors.New("telegram factory is nil")
	}

	if metrics == nil {
		return nil, errors.New("metrics is nil")
	}

	return &notifier{bot: bot, factory: factory, metrics: metrics}, nil
}

// Factory ...
func (tn *notifier) Factory() command2.Factory {
	return tn.factory
}

func (tn *notifier) setKeyboard(message TGBotAPI.MessageConfig) TGBotAPI.MessageConfig {
	message.ReplyMarkup = TGBotAPI.NewReplyKeyboard(
		TGBotAPI.NewKeyboardButtonRow(
			TGBotAPI.NewKeyboardButton(string(command2.ButtonWorkingTime)),
			TGBotAPI.NewKeyboardButton(string(command2.ButtonStatCurrentWorkingPeriod)),
		),
	)

	return message
}

func (tn *notifier) sendMessage(text string, telegramID int64) error {
	_, err := tn.bot.Send(tn.setKeyboard(TGBotAPI.NewMessage(
		telegramID,
		text,
	)))
	if err != nil {
		return fmt.Errorf("error send message: %w", err)
	}

	tn.metrics.IncMessageCounter()

	return nil
}
