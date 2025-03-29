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
}

type pyroscopeConfig struct {
	host string
	port string
}

// NewPyroscopeConfig ...
func NewPyroscopeConfig(os OS) (PyroscopeConfig, error) {
	if os == nil {
		return nil, fmt.Errorf("os must not be nil")
	}

	host := os.Getenv(pyroscopeHostEnvName)
	if len(host) == 0 {
		return nil, fmt.Errorf("environment variable %s not set", pyroscopeHostEnvName)
	}

	port := os.Getenv(pyroscopePortEnvName)
	if len(port) == 0 {
		return nil, fmt.Errorf("environment variable %s not set", pyroscopePortEnvName)
	}

	return &pyroscopeConfig{
		host: host,
		port: port,
	}, nil
}

// Address ...
func (c *pyroscopeConfig) Address() string {
	return net.JoinHostPort(c.host, c.port)
}
