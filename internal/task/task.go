package task

import "context"

// Task ???
type Task interface {
	Run(ctx context.Context) error
	GetName() string
}
