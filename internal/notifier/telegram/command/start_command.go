package command

import (
	"context"
)

// StartCommand Команда /start
type StartCommand struct{}

// GetMessage Метод получения текста сообщения в ответ на команду /start
func (StartCommand) GetMessage(_ context.Context) (string, error) {
	return "Добро пожаловать. Используйте кнопки для получения информации", nil
}
