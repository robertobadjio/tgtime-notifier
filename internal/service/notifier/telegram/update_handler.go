package telegram

import (
	"context"
	"fmt"

	"github.com/robertobadjio/tgtime-notifier/internal/helper"
)

// SendCommandMessage Метод для отправки сообщения в ответ на команду пользователя.
func (tn *TGNotifier) SendCommandMessage(ctx context.Context, p ParamsUpdate) error {
	switch p.Update.Message.Text {
	case ButtonStart:
		return tn.sendMessage(tn.handleStartCommand(), int64(p.Update.Message.From.ID))
	case ButtonWorkingTime:
		message, err := tn.handleWorkingTimeCommand(ctx, int64(p.Update.Message.From.ID))
		if err != nil {
			return fmt.Errorf("handle working time command: %w", err)
		}

		return tn.sendMessage(message, int64(p.Update.Message.From.ID))
	case ButtonStatCurrentWorkingPeriod:
		return tn.sendMessage(tn.handleStatCurrentWorkingPeriodCommand(ctx), int64(p.Update.Message.From.ID))
	}

	return nil
}

func (tn *TGNotifier) handleStartCommand() string {
	return "Добро пожаловать. Используйте кнопки для получения информации."
}

func (tn *TGNotifier) handleStatCurrentWorkingPeriodCommand(_ context.Context) string {
	// TODO: Implement method

	/*
		} else if Type(update.Message.Text) == ButtonStatCurrentWorkingPeriod {
				// TODO: Реализовать в tgtime-api метод получения идентификатора текущего периода
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
					int64(TelegramID),
					fmt.Sprintf(
						"Статистика за период с %s до %s\nВсего в этом месяце %d из %d часов",
						result.StartWorkingDate,
						result.EndWorkingDate,
						result.WorkingHours,
						result.TotalWorkingHours,
					))

				t.bot.Send(t.setKeyboard(message))
			}
	*/
	return ""
}

func (tn *TGNotifier) handleWorkingTimeCommand(ctx context.Context, telegramID int64) (string, error) {
	user, err := tn.TGTimeAPIClient.GetUserByTelegramID(ctx, telegramID)
	if err != nil {
		return "", fmt.Errorf("error getting user by telegram ID: %w", err)
	}

	timeSummaryResponse, err := tn.TGTimeAggregatorClient.GetTimeSummary(
		ctx,
		user.User.MacAddress,
		helper.GetNow().Format("2006-01-02"),
	)
	if err != nil {
		return "", fmt.Errorf("error getting time summary: %w", err)
	}
	if len(timeSummaryResponse.Summary) == 0 {
		return "", fmt.Errorf("time summary not found")
	}

	if timeSummaryResponse.Summary[0].SecondsStart == 0 {
		return "Вы сегодня не были в офисе", nil
	}

	hours, minutes := helper.SecondsToHM(timeSummaryResponse.Summary[0].Seconds)
	beginTime := helper.SecondsToTime(timeSummaryResponse.Summary[0].SecondsStart)
	mes := fmt.Sprintf(
		"Сегодня Вы в офисе с %s\nУчтенное время %d ч. %d м.",
		beginTime.Format("15:04"),
		hours,
		minutes,
	)

	breaks := helper.BreaksToString(helper.BuildBreaks(timeSummaryResponse.Summary[0].GetBreaks()))
	if breaks != "" {
		mes += fmt.Sprintf("\nПерерывы %s", breaks)
	}

	return mes, nil
}
