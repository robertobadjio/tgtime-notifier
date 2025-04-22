package previous_day_info

import (
	"context"
	"fmt"
	"time"

	"github.com/robertobadjio/tgtime-notifier/internal/helper"
	"github.com/robertobadjio/tgtime-notifier/internal/service/client/aggregator"
	"github.com/robertobadjio/tgtime-notifier/internal/service/client/api_pb"
	"github.com/robertobadjio/tgtime-notifier/internal/service/notifier/telegram"
)

// PreviousDayInfo ...
type PreviousDayInfo interface {
	Run(ctx context.Context) error
}

type notifier interface {
	SendPreviousDayInfoMessage(ctx context.Context, params telegram.ParamsPreviousDayInfo) error
}

type previousDayInfo struct {
	tgTimeAPIClient        api_pb.Client
	tgTimeAggregatorClient aggregator.Client
	notifier               notifier
}

// NewPreviousDayInfo ???
func NewPreviousDayInfo(
	tgTimeAPIClient api_pb.Client,
	tgTimeAggregatorClient aggregator.Client,
	notifier notifier,
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
		user, errGetUserByMACAddress := pdi.tgTimeAPIClient.GetUserByMacAddress(ctx, summaryByUser.MacAddress)
		if errGetUserByMACAddress != nil {
			// TODO: log error
			continue
			//return fmt.Errorf("error getting user by telegram id: %w", err)
		}

		hours, minutes := helper.SecondsToHM(summaryByUser.GetSeconds())

		errSendPreviousDayInfo := pdi.notifier.SendPreviousDayInfoMessage(
			ctx,
			telegram.ParamsPreviousDayInfo{
				TelegramID:   user.GetUser().TelegramId,
				SecondsStart: helper.SecondsToTime(summaryByUser.GetSecondsStart()),
				SecondsEnd:   helper.SecondsToTime(summaryByUser.GetSecondsEnd()),
				Hours:        hours,
				Minutes:      minutes,
				Breaks:       helper.BreaksToString(helper.BuildBreaks(summaryByUser.Breaks)),
			},
		)
		if errSendPreviousDayInfo != nil {
			return fmt.Errorf("error sending previous day info message: %w", errSendPreviousDayInfo)
		}
	}

	return nil
}

func getPreviousDate(location string) time.Time {
	moscowLocation, _ := time.LoadLocation(location)
	return time.Now().AddDate(0, 0, -1).In(moscowLocation)
}
