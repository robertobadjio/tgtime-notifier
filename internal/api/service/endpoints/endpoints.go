package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/robertobadjio/tgtime-notifier/internal/api"
)

// Set ...
type Set struct {
	LivenessEndpoint  endpoint.Endpoint
	ReadinessEndpoint endpoint.Endpoint
}

// NewEndpointSet ...
func NewEndpointSet(s *api.NotifierService) Set {
	return Set{
		LivenessEndpoint:  MakeLivenessEndpoint(s),
		ReadinessEndpoint: MakeReadinessEndpoint(s),
	}
}

// MakeLivenessEndpoint ...
func MakeLivenessEndpoint(svc *api.NotifierService) endpoint.Endpoint {
	return func(_ context.Context, _ interface{}) (interface{}, error) {
		return LivenessResponse{Code: svc.Liveness(), Err: ""}, nil
	}
}

// MakeReadinessEndpoint ...
func MakeReadinessEndpoint(svc *api.NotifierService) endpoint.Endpoint {
	return func(_ context.Context, _ interface{}) (interface{}, error) {
		return ReadinessResponse{Code: svc.Readiness(), Err: ""}, nil
	}
}
