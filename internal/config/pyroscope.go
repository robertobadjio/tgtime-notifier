package config

import (
	"fmt"
	"net"
)

const pyroscopeHostEnvName = "PYROSCOPE_HOST"
const pyroscopePortEnvName = "PYROSCOPE_PORT"

const applicationName string = "notify.app"

// PyroscopeConfig ...
type PyroscopeConfig struct {
	host            string
	port            string
	applicationName string
	enabled         bool
}

// NewPyroscopeConfig ...
func NewPyroscopeConfig(os OS) (*PyroscopeConfig, error) {
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

	return &PyroscopeConfig{
		host:            host,
		port:            port,
		applicationName: applicationName,
		enabled:         enabled,
	}, nil
}

// Address ...
func (c *PyroscopeConfig) Address() string {
	return net.JoinHostPort(c.host, c.port)
}

// Enabled ...
func (c *PyroscopeConfig) Enabled() bool {
	return c.enabled
}

// ApplicationName ...
func (c *PyroscopeConfig) ApplicationName() string {
	return c.applicationName
}
