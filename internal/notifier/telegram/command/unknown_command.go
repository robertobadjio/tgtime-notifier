package command

import (
	"context"
)

// UnknownCommand Неизвестная команда
type UnknownCommand struct{}

// GetMessage Метод получения текста сообщения о неизвестной команде
func (UnknownCommand) GetMessage(_ context.Context) (string, error) {
	return "Неизвестная команда", nil
}
