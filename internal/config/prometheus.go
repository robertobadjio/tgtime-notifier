package config

import (
	"fmt"
	"net"
)

const promAppPortEnvName = "PROMETHEUS_APP_PORT"
const endpointPath = "/metrics"

// PromConfig ...
type PromConfig struct {
	host    string
	port    string
	path    string
	enabled bool
}

// NewPromConfig ...
func NewPromConfig(os OS) (*PromConfig, error) {
	if os == nil {
		return nil, fmt.Errorf("os must not be nil")
	}

	port := os.Getenv(promAppPortEnvName)
	enabled := true
	if len(port) == 0 {
		enabled = false
	}

	return &PromConfig{
		host:    "",
		port:    port,
		path:    endpointPath,
		enabled: enabled,
	}, nil
}

// Address ...
func (c *PromConfig) Address() string {
	return net.JoinHostPort(c.host, c.port)
}

// Path ...
func (c *PromConfig) Path() string {
	return c.path
}

// Enabled ...
func (c *PromConfig) Enabled() bool {
	return c.enabled
}
