package config

import (
	"fmt"
	"net"
)

const (
	tgTimeAPIHost = "TGTIME_API_HOST"
	tgTimeAPIPort = "TGTIME_API_PORT"
)

// TgTimeAPIConfig ???
type TgTimeAPIConfig interface {
	Address() string
}

type tgTimeAPIConfig struct {
	host string
	port string
}

// NewTgTimeAPIConfig ???
func NewTgTimeAPIConfig(os OS) (TgTimeAPIConfig, error) {
	if os == nil {
		return nil, fmt.Errorf("os must not be nil")
	}

	host := os.Getenv(tgTimeAPIHost)

	port := os.Getenv(tgTimeAPIPort)
	if len(port) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", tgTimeAPIPort)
	}

	return &tgTimeAPIConfig{
		host: host,
		port: port,
	}, nil
}

// Address ???
func (tta *tgTimeAPIConfig) Address() string {
	return net.JoinHostPort(tta.host, tta.port)
}
