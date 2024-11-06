package previous_day_info

import (
	"context"
	"fmt"
	"time"

	"github.com/robertobadjio/tgtime-notifier/internal/aggregator"
	"github.com/robertobadjio/tgtime-notifier/internal/api_pb"
	"github.com/robertobadjio/tgtime-notifier/internal/notifier"
	"github.com/robertobadjio/tgtime-notifier/internal/notifier/telegram"
	"github.com/robertobadjio/tgtime-notifier/internal/notifier/telegram/command"
	"github.com/robertobadjio/tgtime-notifier/internal/task"
)

const taskName = "previous_day_info"

type previousDayInfoTask struct {
	tgTimeAPIClient        api_pb.Client
	tgTimeAggregatorClient aggregator.Client
	notifier               notifier.Notifier
}

// NewPreviousDayInfoTask ???
func NewPreviousDayInfoTask(
	tgTimeAPIClient api_pb.Client,
	tgTimeAggregatorClient aggregator.Client,
	notifier notifier.Notifier,
) task.Task {
	return &previousDayInfoTask{
		tgTimeAPIClient:        tgTimeAPIClient,
		tgTimeAggregatorClient: tgTimeAggregatorClient,
		notifier:               notifier,
	}
}

// GetName ???
func (t *previousDayInfoTask) GetName() string {
	return taskName
}

func (t *previousDayInfoTask) Run(ctx context.Context) error {
	timeSummaryResponse, err := t.tgTimeAggregatorClient.GetTimeSummary(
		ctx,
		"",
		getPreviousDate("Europe/Moscow").Format("2006-01-02"),
	)
	if err != nil {
		return fmt.Errorf("error getting time summary: %w", err)
	}
	if len(timeSummaryResponse.Summary) == 0 {
		return nil
	}

	for _, summaryByUser := range timeSummaryResponse.Summary {
		user, err := t.tgTimeAPIClient.GetUserByMacAddress(ctx, summaryByUser.MacAddress)
		if err != nil {
			// TODO: log error
			continue
			//return fmt.Errorf("error getting user by telegram id: %w", err)
		}

		hours, minutes := command.SecondsToHM(summaryByUser.GetSeconds())

		err = t.notifier.SendPreviousDayInfoMessage(
			ctx,
			telegram.ParamsPreviousDayInfo{
				TelegramID:   user.GetUser().TelegramId,
				SecondsStart: command.SecondsToTime(summaryByUser.GetSecondsStart()),
				SecondsEnd:   command.SecondsToTime(summaryByUser.GetSecondsEnd()),
				Hours:        hours,
				Minutes:      minutes,
				Breaks:       command.BreaksToString(command.BuildBreaks(summaryByUser.Breaks)),
			},
		)
		if err != nil {
			return fmt.Errorf("error sending previous day info message: %w", err)
		}
	}

	return nil
}

func getPreviousDate(location string) time.Time {
	moscowLocation, _ := time.LoadLocation(location)
	return time.Now().AddDate(0, 0, -1).In(moscowLocation)
}
