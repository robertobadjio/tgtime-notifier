package config

import (
	"fmt"
	"net"
)

const promAppPortEnvName = "PROMETHEUS_APP_PORT"

// PromConfig ...
type PromConfig interface {
	Address() string
	Path() string
}

type promConfig struct {
	host string
	port string
	path string
}

// NewPromConfig ...
func NewPromConfig(os OS) (PromConfig, error) {
	if os == nil {
		return nil, fmt.Errorf("os must not be nil")
	}

	port := os.Getenv(promAppPortEnvName)
	if len(port) == 0 {
		return nil, fmt.Errorf("environment variable %s not set", promAppPortEnvName)
	}

	return &promConfig{
		host: "",
		port: port,
		path: "/metrics",
	}, nil
}

// Address ...
func (c *promConfig) Address() string {
	return net.JoinHostPort(c.host, c.port)
}

// Path ...
func (c *promConfig) Path() string {
	return c.path
}
