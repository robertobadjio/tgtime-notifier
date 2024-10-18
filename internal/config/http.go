package config

import (
	"fmt"
	"net"
	"os"
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
func NewHTTPConfig() (HTTPConfig, error) {
	port := os.Getenv(httpPortEnvVar)
	if len(port) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", httpPortEnvVar)
	}

	return &httpConfig{host: "", port: port}, nil
}

func (h *httpConfig) Address() string {
	return net.JoinHostPort(h.host, h.port)
}
