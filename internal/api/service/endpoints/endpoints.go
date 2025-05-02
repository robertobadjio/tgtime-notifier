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
	return func(_ context.Context, request interface{}) (interface{}, error) {
		_ = request.(LivenessRequest)
		code, err := svc.Liveness()
		if err != nil {
			return LivenessResponse{Code: code, Err: err.Error()}, err
		}
		return LivenessResponse{Code: code, Err: ""}, nil
	}
}

// MakeReadinessEndpoint ...
func MakeReadinessEndpoint(svc *api.NotifierService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		_ = request.(ReadinessRequest)
		code, err := svc.Readiness()
		if err != nil {
			return ReadinessResponse{Code: code, Err: err.Error()}, err
		}
		return ReadinessResponse{Code: code, Err: ""}, nil
	}
}
