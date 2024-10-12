package telegram

import (
	"context"
)

type Context interface {
	SendMessage(ctx context.Context) (string, error)
}

type TypeMessage struct {
	Message Context
}

func (tm *TypeMessage) Handle(ctx context.Context) (string, error) {
	return tm.Message.SendMessage(ctx)
}
