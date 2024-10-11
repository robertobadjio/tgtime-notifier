package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"tgtime-notifier/internal/aggregator"
	"tgtime-notifier/internal/api_pb"
	"tgtime-notifier/internal/notifier"

	//"tgtime-notifier/internal/api"
	//"tgtime-notifier/internal/notifier"
	"time"
)

func (t *TelegramNotifier) Info(
	ctx context.Context,
	update tgbotapi.Update,
	clientAggregator *aggregator.Client,
	clientApi *api_pb.Client,
) error {
	telegramId := update.Message.From.ID

	// TODO: Стратегия
	if update.Message.Text == buttonWorkingTime {
		// TODO: Получить mac address по tg id
		user, err := clientApi.GetUserByTelegramId(ctx, int64(telegramId))
		if err != nil {
			return fmt.Errorf("error getting user by telegram id: %w", err)
		}

		timeSummary, err := clientAggregator.GetTimeSummary(
			ctx,
			user.User.MacAddress,
			getNow().Format("2006-01-02"),
		)
		if err != nil {
			return fmt.Errorf("error getting time summary: %w", err)
		}

		var messageTelegram tgbotapi.MessageConfig
		if timeSummary.TimeSummary[0].SecondsStart == 0 {
			messageTelegram = tgbotapi.NewMessage(int64(telegramId), "Вы сегодня не были в офисе")
		} else {
			hours, minutes := secondsToHM(int(timeSummary.TimeSummary[0].Seconds))
			beginTime := time.Unix(timeSummary.TimeSummary[0].SecondsStart, 0)
			message := fmt.Sprintf(
				"Сегодня Вы в офисе с %s\nУчтенное время %d ч. %d м.",
				beginTime.Format("15:04"),
				hours,
				minutes,
			)

			var breaksRaw []*notifier.Break

			// TODO: По GRPC отдавать сразу срез
			_ = json.Unmarshal([]byte(timeSummary.TimeSummary[0].GetBreaksJson()), &breaksRaw)
			breaks := breaksToString(buildBreaks(breaksRaw))
			if breaks != "" {
				message += fmt.Sprintf("\nПерерывы %s", breaks)
			}
			messageTelegram = tgbotapi.NewMessage(int64(telegramId), message)
		}
		_, err = t.bot.Send(t.setKeyboard(messageTelegram))
		if err != nil {
			return fmt.Errorf("error sending telegram message - working time: %w", err)
		}
	} else if update.Message.Text == buttonStatCurrentWorkingPeriod {
		// TODO: Реализовать в tgtime-api метод полчения идентификатора текущего периода
		// Период и даты получили.
		// TODO: С переодом нужно получить result.TotalWorkingHours, - рабочее колисчество часов в периоде
		// StartWorkingDate, err := time.Parse(time.RFC3339, period.BeginDate)
		//	if err != nil {
		//		panic(err)
		//	}
		//	EndWorkingDate, err := time.Parse(time.RFC3339, period.EndDate)
		//	if err != nil {
		//		panic(err)
		//	}
		// StartWorkingDate:  start.Format("02.01.2006"),
		//		EndWorkingDate:    end.Format("02.01.2006"),
		// TODO: Запрашиваем в tgtime-aggregator time summary по macAddress и dates
		// var totalMonthWorkingTime int64
		//	for _, timeResponse := range periodUser.Time {
		//		totalMonthWorkingTime += timeResponse.Total
		//	}
		// WorkingHours:      totalMonthWorkingTime / 3600,

		/*message := tgbotapi.NewMessage(
			int64(telegramId),
			fmt.Sprintf(
				"Статистика за период с %s до %s\nВсего в этом месяце %d из %d часов",
				result.StartWorkingDate,
				result.EndWorkingDate,
				result.WorkingHours,
				result.TotalWorkingHours,
			))

		t.bot.Send(t.setKeyboard(message))*/
	} else if update.Message.Text == buttonStart {
		message := tgbotapi.NewMessage(
			int64(telegramId),
			"Добро пожаловать. Используйте кнопки для получения информации")
		_, err := t.bot.Send(t.setKeyboard(message))
		if err != nil {
			return fmt.Errorf("error sending telegram message - welcome message: %w", err)
		}
	} else {
		message := tgbotapi.NewMessage(
			int64(telegramId),
			"Неизвестная команда")
		_, err := t.bot.Send(t.setKeyboard(message))
		if err != nil {
			return fmt.Errorf("error sending telegram message - unknown command: %w", err)
		}
	}

	return nil
}

func secondsToHM(seconds int) (int, int) {
	hours := seconds / 3600
	minutes := (seconds / 60) - (hours * 60)

	return hours, minutes
}

func getNow() time.Time {
	return time.Now().In(getMoscowLocation())
}

func getMoscowLocation() *time.Location {
	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	return moscowLocation
}
