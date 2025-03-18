package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/pprof"

	httpTransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/robertobadjio/tgtime-notifier/internal/api/service/endpoints"
)

type errorer interface {
	Error() error
}

// HandlerWithError ...
type HandlerWithError func(w http.ResponseWriter, r *http.Request) error

// NewHTTPHandler ...
func NewHTTPHandler(ep endpoints.Set) http.Handler {
	router := mux.NewRouter()

	opt := []httpTransport.ServerOption{
		httpTransport.ServerErrorEncoder(encodeError),
	}

	router.Methods(http.MethodGet).Path("/liveness").Handler(httpTransport.NewServer(
		ep.LivenessEndpoint,
		decodeHTTPLivenessRequest,
		encodeResponse,
		opt...,
	))

	router.Methods(http.MethodGet).Path("/readiness").Handler(httpTransport.NewServer(
		ep.ReadinessEndpoint,
		decodeHTTPReadinessRequest,
		encodeResponse,
		opt...,
	))

	router.Methods(http.MethodGet).Path("/debug/pprof/profile").HandlerFunc(pprof.Profile)
	router.Methods(http.MethodGet).Path("/debug/pprof/trace").HandlerFunc(pprof.Trace)
	router.Methods(http.MethodGet).Path("/debug/pprof/heap").Handler(pprof.Handler("heap"))

	return router
}

func decodeHTTPLivenessRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	var req endpoints.LivenessRequest
	return req, nil
}

func decodeHTTPReadinessRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	var req endpoints.ReadinessRequest
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.Error() != nil {
		encodeError(ctx, e.Error(), w)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{ // TODO: Handle error
		"error": err.Error(),
	})
}
