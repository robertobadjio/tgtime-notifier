package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (t *TelegramNotifier) SendWelcomeMessage(_ context.Context, telegramId int64) error {
	msg := tgbotapi.NewMessage(telegramId, "Вы пришли в офис")
	_, err := t.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("sending welcome message failed: %w", err)
	}

	return nil
}
