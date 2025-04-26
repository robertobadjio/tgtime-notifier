package previous_day_info

import (
	"context"
	"fmt"
	"time"

	"github.com/jonboulle/clockwork"

	pb "github.com/robertobadjio/tgtime-aggregator/pkg/api/time_v1"
	pbapiv1 "github.com/robertobadjio/tgtime-api/api/v1/pb/api"
	"github.com/robertobadjio/tgtime-notifier/internal/helper"
	"github.com/robertobadjio/tgtime-notifier/internal/logger"
	"github.com/robertobadjio/tgtime-notifier/internal/service/notifier/telegram"
)

type notifier interface {
	SendPreviousDayInfoMessage(ctx context.Context, params telegram.ParamsPreviousDayInfo) error
}

type aggregatorClient interface {
	GetTimeSummary(
		ctx context.Context,
		macAddress, date string,
	) (*pb.GetSummaryResponse, error)
}

// APIClient ...
type APIClient interface {
	GetUserByMacAddress(
		ctx context.Context,
		macAddress string,
	) (*pbapiv1.GetUserByMacAddressResponse, error)
}

// PreviousDayInfo ...
type PreviousDayInfo struct {
	APIClient              APIClient
	tgTimeAggregatorClient aggregatorClient
	notifier               notifier
	clock                  clockwork.Clock
	hour, minute, second   int
}

// NewPreviousDayInfo ...
func NewPreviousDayInfo(
	tgTimeAPIClient APIClient,
	aggregatorClient aggregatorClient,
	notifier notifier,
	c clockwork.Clock,
	hour, minute, second int,
) (*PreviousDayInfo, error) {
	if c == nil {
		return nil, fmt.Errorf("click must be set")
	}

	if hour < 0 || hour > 23 {
		return nil, fmt.Errorf("hour must be between 0 and 23")
	}

	if minute < 0 || minute > 59 {
		return nil, fmt.Errorf("minute must be between 0 and 59")
	}

	if second < 0 || second > 59 {
		return nil, fmt.Errorf("second must be between 0 and 59")
	}

	if nil == tgTimeAPIClient {
		return nil, fmt.Errorf("APIClient must be set")
	}

	if nil == aggregatorClient {
		return nil, fmt.Errorf("aggregatorClient must be set")
	}

	if nil == notifier {
		return nil, fmt.Errorf("notifier must be set")
	}

	return &PreviousDayInfo{
		APIClient:              tgTimeAPIClient,
		tgTimeAggregatorClient: aggregatorClient,
		notifier:               notifier,
		clock:                  c,
		hour:                   hour,
		minute:                 minute,
		second:                 second,
	}, nil
}

// Start ...
func (pdi *PreviousDayInfo) Start(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	go func() {
		_ = pdi.everyDayByHour(ctx, pdi.sendAllUsersNotify, time.NewTicker(time.Minute))
	}()

	return nil
}

func (pdi *PreviousDayInfo) everyDayByHour(
	ctx context.Context,
	handler func(ctx context.Context, date string) error,
	ticker *time.Ticker,
) error {
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		select {
		case <-ticker.C:
			cc := pdi.clock
			h, m, s := cc.Now().Clock()
			if h == pdi.hour && m == pdi.minute && s == pdi.second {
				err := handler(ctx, pdi.getPreviousDate("Europe/Moscow").Format(time.DateOnly))
				if err != nil {
					logger.Error(
						"component", "previous_day_info",
						"during", "send notify",
						"error", err.Error(),
					)
				}
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// Run ...
func (pdi *PreviousDayInfo) sendAllUsersNotify(ctx context.Context, date string) error {
	timeSummaryResponse, err := pdi.tgTimeAggregatorClient.GetTimeSummary(ctx, "", date)
	if err != nil {
		return fmt.Errorf("error getting time summary: %w", err)
	}

	if len(timeSummaryResponse.Summary) == 0 {
		return nil
	}

	for _, summaryByUser := range timeSummaryResponse.Summary {
		user, errGetUserByMACAddress := pdi.APIClient.GetUserByMacAddress(ctx, summaryByUser.MacAddress)
		if errGetUserByMACAddress != nil {
			logger.Warn(
				"component", "previous_day_info",
				"during", "send notify",
				"error", errGetUserByMACAddress.Error(),
			)
			continue
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

func (pdi *PreviousDayInfo) getPreviousDate(location string) time.Time {
	moscowLocation, _ := time.LoadLocation(location)
	return time.Now().AddDate(0, 0, -1).In(moscowLocation)
}
