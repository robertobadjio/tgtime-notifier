package config

import (
	"fmt"
	"net"
)

const httpPortEnvVar = "HTTP_PORT"

// HTTPConfig ???
type HTTPConfig interface {
	Address() string
}

type httpConfig struct {
	host string
	port string
}

// NewHTTPConfig ???
func NewHTTPConfig(os OS) (HTTPConfig, error) {
	if os == nil {
		return nil, fmt.Errorf("os must not be nil")
	}

	port := os.GetEnv(httpPortEnvVar)
	if len(port) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", httpPortEnvVar)
	}

	return &httpConfig{host: "", port: port}, nil
}

func (h *httpConfig) Address() string {
	return net.JoinHostPort(h.host, h.port)
}
