package telegram

import (
	"context"
	"fmt"

	notifierI "github.com/robertobadjio/tgtime-notifier/internal/service/notifier"
)

// SendCommandMessage Метод для отправки сообщения в ответ на команду пользователя
func (tn *notifier) SendCommandMessage(ctx context.Context, params notifierI.Params) error {
	p, ok := params.(ParamsUpdate)
	if !ok {
		return fmt.Errorf("error cast interface param")
	}

	message, err := tn.Factory().GetCommandHandler(p.Update).GetMessage(ctx)
	if err != nil {
		return fmt.Errorf("error getting text message: %w", err)
	}

	return tn.SendMessage(message, int64(p.Update.Message.From.ID))
}
