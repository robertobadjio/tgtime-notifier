package telegram

import (
	"context"
	"fmt"

	notifierI "github.com/robertobadjio/tgtime-notifier/internal/notifier"
)

// SendWelcomeMessage ???
func (tn notifier) SendWelcomeMessage(_ context.Context, params notifierI.Params) error {
	p, ok := params.(ParamsWelcomeMessage)
	if !ok {
		return fmt.Errorf("error cast interface param")
	}

	return tn.SendMessage("Вы пришли в офис", p.TelegramID)
}
