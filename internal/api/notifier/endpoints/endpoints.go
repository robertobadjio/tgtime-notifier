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
func NewEndpointSet(s api.Service) Set {
	return Set{
		LivenessEndpoint:  MakeLivenessEndpoint(s),
		ReadinessEndpoint: MakeReadinessEndpoint(s),
	}
}

// MakeLivenessEndpoint ...
func MakeLivenessEndpoint(svc api.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(LivenessRequest)
		code, err := svc.Liveness(ctx)
		if err != nil {
			return LivenessResponse{Code: code, Err: err.Error()}, err
		}
		return LivenessResponse{Code: code, Err: ""}, nil
	}
}

// MakeReadinessEndpoint ...
func MakeReadinessEndpoint(svc api.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(ReadinessRequest)
		code, err := svc.Readiness(ctx)
		if err != nil {
			return ReadinessResponse{Code: code, Err: err.Error()}, err
		}
		return ReadinessResponse{Code: code, Err: ""}, nil
	}
}
