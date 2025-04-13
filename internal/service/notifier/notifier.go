package notifier

import (
	"context"
)

// Params ...
type Params interface{}

// Notifier ...
type Notifier interface {
	SendCommandMessage(ctx context.Context, params Params) error
	SendWelcomeMessage(ctx context.Context, params Params) error
	SendPreviousDayInfoMessage(ctx context.Context, params Params) error
}
