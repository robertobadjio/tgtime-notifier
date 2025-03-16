package api

import (
	"context"
	"net/http"
)

func (s *notifierService) Readiness(_ context.Context) (int, error) {
	return http.StatusOK, nil
}
