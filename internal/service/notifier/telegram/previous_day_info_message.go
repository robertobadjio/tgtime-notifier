package telegram

import (
	"context"
	"fmt"
)

// SendPreviousDayInfoMessage ...
func (tn *TGNotifier) SendPreviousDayInfoMessage(_ context.Context, p ParamsPreviousDayInfo) error {
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

	return tn.sendMessage(message, p.TelegramID)
}
