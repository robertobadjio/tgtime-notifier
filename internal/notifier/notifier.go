package notifier

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Notifier interface {
	SendWelcomeMessage(ctx context.Context, telegramId int64)
	SendMessageCommand(ctx context.Context, update tgbotapi.Update)
}
