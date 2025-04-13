package command

import (
	"context"
)

type startCommand struct{}

// NewStartCommand ...
func NewStartCommand() Command {
	return &startCommand{}
}

// GetMessage Метод получения текста сообщения в ответ на команду /start.
func (startCommand) GetMessage(_ context.Context) (string, error) {
	return "Добро пожаловать. Используйте кнопки для получения информации", nil
}
