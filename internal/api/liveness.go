package api

import (
	"net/http"
)

// Liveness ...
func (s *NotifierService) Liveness() int {
	return http.StatusOK
}
