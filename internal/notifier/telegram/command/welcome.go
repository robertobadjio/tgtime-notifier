package command

import (
	"context"
)

// WelcomeCommand Сообщение по приходе сотрудника в офис / на работу
type WelcomeCommand struct{}

// GetMessage Метод получения текста сообщения по приходе сотрудника в офис / на работу
func (WelcomeCommand) GetMessage(_ context.Context) (string, error) {
	return "Вы пришли в офис", nil
}
