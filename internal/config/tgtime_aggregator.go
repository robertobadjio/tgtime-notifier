package config

import (
	"fmt"
	"net"
	"os"
)

const (
	tgTimeAggregatorHost = "TGTIME_AGGREGATOR_HOST"
	tgTimeAggregatorPort = "TGTIME_AGGREGATOR_PORT"
)

// TgTimeAggregatorConfig ???
type TgTimeAggregatorConfig interface {
	Address() string
}

type tgTimeAggregatorConfig struct {
	host string
	port string
}

// NewTgTimeAggregatorConfig ???
func NewTgTimeAggregatorConfig() (TgTimeAggregatorConfig, error) {
	host := os.Getenv(tgTimeAggregatorHost)

	port := os.Getenv(tgTimeAggregatorPort)
	if len(port) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", tgTimeAggregatorPort)
	}

	return &tgTimeAggregatorConfig{
		host: host,
		port: port,
	}, nil
}

// Address ???
func (tta *tgTimeAggregatorConfig) Address() string {
	return net.JoinHostPort(tta.host, tta.port)
}
