package telegram

import (
	"context"

	"github.com/robertobadjio/tgtime-notifier/internal/notifier/telegram/command"
)

// Command Интерфейс команды пользователя
type Command interface {
	GetMessage(ctx context.Context) (string, error)
}

// MessageType Тип сообщения - команда пользователя
type MessageType string

const (
	buttonWorkingTime              MessageType = "⏳ Рабочее время"
	buttonStatCurrentWorkingPeriod MessageType = "🗓 Статистика за рабочий период"
	buttonStart                    MessageType = "/start"
	welcome                        MessageType = "welcome"
	//previousDayInfo                        MessageType = "previousDayInfo"
)

// NewCommand Фабрика для получения обработчика команды пользователя
func NewCommand(t MessageType, telegramID int64) Command {
	switch t {
	case buttonStart:
		return command.StartCommand{}
	case buttonWorkingTime:
		return command.WorkingTimeCommand{TelegramID: telegramID}
	case welcome:
		return command.WelcomeCommand{}
	default:
		return command.UnknownCommand{}
	}
}
