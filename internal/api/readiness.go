package api

import (
	"net/http"
)

// Readiness ...
func (s *NotifierService) Readiness() (int, error) {
	return http.StatusOK, nil
}
