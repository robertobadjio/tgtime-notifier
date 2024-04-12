package service

import (
	"cloud-time-tracker/cmd/officetime/api"
	"time"
)

func GetCurrentPeriod() int {
	c := api.NewOfficeTimeClient()
	periods, err := c.GetAllPeriods()
	if err != nil {
		panic(err)
	}

	for _, period := range periods.Periods {
		// TODO: Начало и окончания каждого месяца мне входят в интервал
		// if GetNow().After(GetTimeFromStringDate(period.BeginDate)) && GetNow().Before(GetTimeFromStringDate(period.EndDate).Add(time.Hour * 24)) {
		if GetNow().After(GetTimeFromStringDate(period.BeginDate)) && GetNow().Before(GetTimeFromStringDate(period.EndDate)) {
			return period.Id
		}
	}

	return 0
}

func GetTimeFromStringDate(date string) time.Time {
	timeObject, _ := time.ParseInLocation("2006-01-02", date, GetMoscowLocation())

	return timeObject
}

func GetNow() time.Time {
	return time.Now().In(GetMoscowLocation())
}

func GetMoscowLocation() *time.Location {
	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	return moscowLocation
}
