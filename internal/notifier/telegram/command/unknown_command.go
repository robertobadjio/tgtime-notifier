package command

import (
	"context"
)

type UnknownCommand struct{}

func (UnknownCommand) GetMessage(_ context.Context) (string, error) {
	return "Неизвестная команда", nil
}
