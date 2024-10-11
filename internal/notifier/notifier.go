package notifier

import (
	"context"
	"time"
)

type Break struct {
	BeginTime int64 `json:"beginTime"` // TODO: rename StartTime
	EndTime   int64 `json:"endTime"`   // TODO: rename EndTime
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
