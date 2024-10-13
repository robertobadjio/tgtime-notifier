package command

/*
} else if MessageType(update.Message.Text) == buttonStatCurrentWorkingPeriod {
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
			int64(telegramId),
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
