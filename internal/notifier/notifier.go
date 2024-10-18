package notifier

import "context"

// Notifier ???
type Notifier interface {
	SendMessageCommand(ctx context.Context) error
	SendWelcomeMessage(ctx context.Context) error
}
