package command

import (
	"context"
)

type StartCommand struct{}

func (StartCommand) SendMessage(_ context.Context) (string, error) {
	return "Добро пожаловать. Используйте кнопки для получения информации", nil
}
