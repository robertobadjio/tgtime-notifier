package telegram

/*
import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"tgtime-notifier/internal/notifier"
	"time"
)

func (t *TelegramNotifier) SendPreviousDayInfo(
	_ context.Context,
	telegramId int64,
	startTime, endTime time.Time,
	hours, minutes int,
	breaks []*notifier.Break,
) error {
	breaksString := breaksToString(buildBreaks(breaks))
	message := fmt.Sprintf(
		"Вчера Вы были в офисе с %s до %s\nУчтенное время %d ч. %d м.",
		startTime.Format("15:04"),
		endTime.Format("15:04"),
		hours,
		minutes,
	)
	if "" != breaksString {
		message += fmt.Sprintf("\nПерерывы %s\n", breaksString)
	}
	t.bot.Send(t.setKeyboard(tgbotapi.NewMessage(telegramId, message)))
	return nil
}*/
