package command

import (
	"context"
)

type unknownCommand struct{}

// NewUnknownCommand ...
func NewUnknownCommand() Command {
	return &unknownCommand{}
}

// GetMessage Метод получения текста сообщения о неизвестной команде.
func (unknownCommand) GetMessage(_ context.Context) (string, error) {
	return "Неизвестная команда", nil
}
