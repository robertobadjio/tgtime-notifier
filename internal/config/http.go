package config

import (
	"fmt"
	"net"
)

const httpPortEnvVar = "HTTP_PORT"

// HTTPConfig ...
type HTTPConfig struct {
	host string
	port string
}

// NewHTTPConfig ???
func NewHTTPConfig(os OS) (*HTTPConfig, error) {
	if os == nil {
		return nil, fmt.Errorf("os must not be nil")
	}

	port := os.Getenv(httpPortEnvVar)
	if len(port) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", httpPortEnvVar)
	}

	return &HTTPConfig{host: "", port: port}, nil
}

// Address ...
func (h *HTTPConfig) Address() string {
	return net.JoinHostPort(h.host, h.port)
}
