package command

import (
	"context"
)

type UnknownCommand struct{}

func (UnknownCommand) SendMessage(_ context.Context) (string, error) {
	return "Неизвестная команда", nil
}
