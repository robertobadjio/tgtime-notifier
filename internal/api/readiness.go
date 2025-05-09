package api

import (
	"net/http"
)

// Readiness ...
func (s *NotifierService) Readiness() int {
	return http.StatusOK
}
