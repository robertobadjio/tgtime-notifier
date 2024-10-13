package telegram

import (
	"context"
)

type Context interface {
	GetMessage(ctx context.Context) (string, error)
}

type TypeMessage struct {
	Message Context
}

func (tm *TypeMessage) Handle(ctx context.Context) (string, error) {
	return tm.Message.GetMessage(ctx)
}
