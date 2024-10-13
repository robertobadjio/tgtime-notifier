package telegram

import (
	"io"
	"tgtime-notifier/internal/notifier/telegram/command"
)

type MessageTypes interface {
	MessageType(string) (io.ReadWriteCloser, error)
}

type MessageType string

const (
	buttonWorkingTime              MessageType = "⏳ Рабочее время"
	buttonStatCurrentWorkingPeriod MessageType = "🗓 Статистика за рабочий период"
	buttonStart                    MessageType = "/start"
)

const welcome MessageType = "welcome"

func NewCommand(t MessageType, telegramId int64) Context {
	switch t {
	case buttonStart:
		return command.StartCommand{}
	case buttonWorkingTime:
		return command.WorkingTimeCommand{TelegramId: telegramId}
	case welcome:
		return command.WelcomeCommand{}
	default:
		return command.UnknownCommand{}
	}
}
