package previous_day_info

import (
	"context"
	"fmt"
	"time"

	"github.com/robertobadjio/tgtime-notifier/internal/service/client/aggregator"
	"github.com/robertobadjio/tgtime-notifier/internal/service/client/api_pb"
	"github.com/robertobadjio/tgtime-notifier/internal/service/notifier"
	"github.com/robertobadjio/tgtime-notifier/internal/service/notifier/telegram"
	"github.com/robertobadjio/tgtime-notifier/internal/service/notifier/telegram/command"
)

// PreviousDayInfo ...
type PreviousDayInfo interface {
	Run(ctx context.Context) error
}

type previousDayInfo struct {
	tgTimeAPIClient        api_pb.Client
	tgTimeAggregatorClient aggregator.Client
	notifier               notifier.Notifier
}

// NewPreviousDayInfo ???
func NewPreviousDayInfo(
	tgTimeAPIClient api_pb.Client,
	tgTimeAggregatorClient aggregator.Client,
	notifier notifier.Notifier,
) PreviousDayInfo {
	return &previousDayInfo{
		tgTimeAPIClient:        tgTimeAPIClient,
		tgTimeAggregatorClient: tgTimeAggregatorClient,
		notifier:               notifier,
	}
}

// Run ...
func (pdi *previousDayInfo) Run(ctx context.Context) error {
	timeSummaryResponse, err := pdi.tgTimeAggregatorClient.GetTimeSummary(
		ctx,
		"",
		getPreviousDate("Europe/Moscow").Format(time.DateOnly),
	)
	if err != nil {
		return fmt.Errorf("error getting time summary: %w", err)
	}
	if len(timeSummaryResponse.Summary) == 0 {
		return nil
	}

	for _, summaryByUser := range timeSummaryResponse.Summary {
		user, err := pdi.tgTimeAPIClient.GetUserByMacAddress(ctx, summaryByUser.MacAddress)
		if err != nil {
			// TODO: log error
			continue
			//return fmt.Errorf("error getting user by telegram id: %w", err)
		}

		hours, minutes := command.SecondsToHM(summaryByUser.GetSeconds())

		err = pdi.notifier.SendPreviousDayInfoMessage(
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
