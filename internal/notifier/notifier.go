package notifier

import (
	"context"
	"time"
)

type Break struct {
	StartTime int64
	EndTime   int64
}

type Notifier interface {
	SendWelcomeMessage(ctx context.Context, telegramId int64) error
	SendPreviousDayInfo(
		ctx context.Context,
		telegramId int64,
		startTime, endTime time.Time,
		hours, minutes int,
		breaks []*Break,
	) error
}
