package api

import "context"

// Service ...
type Service interface {
	Liveness(ctx context.Context) (int, error)
	Readiness(ctx context.Context) (int, error)
}
