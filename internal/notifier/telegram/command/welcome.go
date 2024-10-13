package command

import (
	"context"
)

type WelcomeCommand struct{}

func (WelcomeCommand) GetMessage(_ context.Context) (string, error) {
	return "Вы пришли в офис", nil
}
