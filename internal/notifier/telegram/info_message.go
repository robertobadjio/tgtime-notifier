package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"tgtime-notifier/internal/api"
	"tgtime-notifier/internal/notifier"
	"time"
)

func (t *TelegramNotifier) Info(_ context.Context, update tgbotapi.Update) {
	telegramId := update.Message.From.ID

	// TODO: стратегия
	if update.Message.Text == buttonWorkingTime {
		c := api.NewOfficeTimeClient()
		response, err := c.GetTimesTelegramIdByDate(telegramId, time.Now())
		if err != nil {
			panic(err)
		}

		var messageTelegram tgbotapi.MessageConfig
		if response.BeginTime == 0 {
			messageTelegram = tgbotapi.NewMessage(int64(telegramId), "Вы сегодня не были в офисе")
		} else {
			hours, minutes := secondsToHM(response.Total)
			beginTime := time.Unix(response.BeginTime, 0)
			message := fmt.Sprintf(
				"Сегодня Вы в офисе с %s\nУчтенное время %d ч. %d м.",
				beginTime.Format("15:04"),
				hours,
				minutes,
			)

			var breaksRaw []*notifier.Break
			for _, br := range response.Breaks {
				breaksRaw = append(breaksRaw, &notifier.Break{
					StartTime: br.BeginTime,
					EndTime:   br.EndTime,
				})
			}

			breaks := breaksToString(buildBreaks(breaksRaw))
			if breaks != "" {
				message += fmt.Sprintf("\nПерерывы %s", breaks)
			}
			messageTelegram = tgbotapi.NewMessage(int64(telegramId), message)
		}
		t.bot.Send(t.setKeyboard(messageTelegram))
	} else if update.Message.Text == buttonStatCurrentWorkingPeriod {
		c := api.NewOfficeTimeClient()
		result, err := c.GetStatByWorkingPeriod(telegramId, getCurrentPeriod())
		if err != nil {
			panic(err)
		}

		message := tgbotapi.NewMessage(
			int64(telegramId),
			fmt.Sprintf(
				"Статистика за период с %s до %s\nВсего в этом месяце %d из %d часов",
				result.StartWorkingDate,
				result.EndWorkingDate,
				result.WorkingHours,
				result.TotalWorkingHours,
			))

		t.bot.Send(t.setKeyboard(message))
	} else if update.Message.Text == buttonStart {
		message := tgbotapi.NewMessage(
			int64(telegramId),
			"Добро пожаловать. Используйте кнопки для получения информации")
		t.bot.Send(t.setKeyboard(message))
	} else {
		message := tgbotapi.NewMessage(
			int64(telegramId),
			"Неизвестная команда")
		t.bot.Send(t.setKeyboard(message))
	}
}

func secondsToHM(seconds int) (int, int) {
	hours := seconds / 3600
	minutes := (seconds / 60) - (hours * 60)

	return hours, minutes
}

func getCurrentPeriod() int {
	c := api.NewOfficeTimeClient()
	periods, err := c.GetAllPeriods()
	if err != nil {
		panic(err)
	}

	for _, period := range periods.Periods {
		// TODO: Начало и окончания каждого месяца мне входят в интервал
		// if GetNow().After(GetTimeFromStringDate(period.BeginDate)) && GetNow().Before(GetTimeFromStringDate(period.EndDate).Add(time.Hour * 24)) {
		if getNow().After(getTimeFromStringDate(period.BeginDate)) && getNow().Before(getTimeFromStringDate(period.EndDate)) {
			return period.Id
		}
	}

	return 0
}

func getTimeFromStringDate(date string) time.Time {
	timeObject, _ := time.ParseInLocation("2006-01-02", date, getMoscowLocation())

	return timeObject
}

func getNow() time.Time {
	return time.Now().In(getMoscowLocation())
}

func getMoscowLocation() *time.Location {
	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	return moscowLocation
}
