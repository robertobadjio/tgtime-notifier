package command

import (
	"context"
)

const (
	// ButtonWorkingTime ...
	ButtonWorkingTime Type = "⏳ Рабочее время"
	// ButtonStatCurrentWorkingPeriod ...
	ButtonStatCurrentWorkingPeriod Type = "🗓 Статистика за рабочий период"
	// ButtonStart ...
	ButtonStart Type = "/start"
)

// Type Тип сообщения - команда пользователя.
type Type string

// Command ...
type Command interface {
	GetMessage(ctx context.Context) (string, error)
}
