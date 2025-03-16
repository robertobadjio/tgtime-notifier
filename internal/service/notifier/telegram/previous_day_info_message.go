package telegram

import (
	"context"
	"fmt"

	notifierI "github.com/robertobadjio/tgtime-notifier/internal/service/notifier"
)

// SendPreviousDayInfoMessage ???
func (tn *notifier) SendPreviousDayInfoMessage(_ context.Context, params notifierI.Params) error {
	p, ok := params.(ParamsPreviousDayInfo)
	if !ok {
		return fmt.Errorf("error cast interface param")
	}

	message := fmt.Sprintf(
		"Вчера Вы были в офисе с %s до %s\nУчтенное время %d ч. %d м.",
		p.SecondsStart.Format("15:04"),
		p.SecondsEnd.Format("15:04"),
		p.Hours,
		p.Minutes,
	)

	if p.Breaks != "" {
		message += fmt.Sprintf("\nПерерывы %s\n", p.Breaks)
	}

	return tn.SendMessage(message, p.TelegramID)
}
