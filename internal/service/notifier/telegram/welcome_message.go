package telegram

import (
	"context"
)

// SendWelcomeMessage ...
func (tn *TGNotifier) SendWelcomeMessage(_ context.Context, p ParamsWelcomeMessage) error {
	return tn.sendMessage("Вы пришли в офис", p.TelegramID)
}
