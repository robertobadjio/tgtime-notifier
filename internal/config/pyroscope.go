package config

import (
	"fmt"
	"net"
)

const pyroscopeHostEnvName = "PYROSCOPE_HOST"
const pyroscopePortEnvName = "PYROSCOPE_PORT"

// PyroscopeConfig ...
type PyroscopeConfig interface {
	Address() string
	Enabled() bool
}

type pyroscopeConfig struct {
	host    string
	port    string
	enabled bool
}

// NewPyroscopeConfig ...
func NewPyroscopeConfig(os OS) (PyroscopeConfig, error) {
	if os == nil {
		return nil, fmt.Errorf("os must not be nil")
	}

	host := os.Getenv(pyroscopeHostEnvName)
	port := os.Getenv(pyroscopePortEnvName)

	enabled := true
	if len(host) == 0 || len(port) == 0 {
		enabled = false
	} else {
		if len(host) == 0 {
			return nil, fmt.Errorf("environment variable %s not set", pyroscopeHostEnvName)
		}

		if len(port) == 0 {
			return nil, fmt.Errorf("environment variable %s not set", pyroscopePortEnvName)
		}
	}

	return &pyroscopeConfig{
		host:    host,
		port:    port,
		enabled: enabled,
	}, nil
}

// Address ...
func (c *pyroscopeConfig) Address() string {
	return net.JoinHostPort(c.host, c.port)
}

// Enabled ...
func (c *pyroscopeConfig) Enabled() bool {
	return c.enabled
}
