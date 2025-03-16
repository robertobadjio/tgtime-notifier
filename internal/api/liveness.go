package api

import (
	"context"
	"net/http"
)

func (s *notifierService) Liveness(_ context.Context) (int, error) {
	return http.StatusOK, nil
}
