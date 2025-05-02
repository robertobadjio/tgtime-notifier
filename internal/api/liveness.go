package api

import (
	"net/http"
)

// Liveness ...
func (s *NotifierService) Liveness() (int, error) {
	return http.StatusOK, nil
}
