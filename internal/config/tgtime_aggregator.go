package config

import (
	"fmt"
	"net"
)

const (
	tgTimeAggregatorHost = "TGTIME_AGGREGATOR_HOST"
	tgTimeAggregatorPort = "TGTIME_AGGREGATOR_PORT"
)

// TgTimeAggregatorConfig ...
type TgTimeAggregatorConfig struct {
	host string
	port string
}

// NewTgTimeAggregatorConfig ...
func NewTgTimeAggregatorConfig(os OS) (*TgTimeAggregatorConfig, error) {
	if os == nil {
		return nil, fmt.Errorf("os must not be nil")
	}

	host := os.Getenv(tgTimeAggregatorHost)

	port := os.Getenv(tgTimeAggregatorPort)
	if len(port) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", tgTimeAggregatorPort)
	}

	return &TgTimeAggregatorConfig{
		host: host,
		port: port,
	}, nil
}

// Address ...
func (tta TgTimeAggregatorConfig) Address() string {
	return net.JoinHostPort(tta.host, tta.port)
}
