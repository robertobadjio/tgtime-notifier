package config

import (
	"fmt"
	"strconv"
)

const (
	kafkaHostEnvName = "KAFKA_BROKER_"
)

// KafkaConfig Конфиг для подключения к Kafka.
type KafkaConfig struct {
	addresses []string
	enabled   bool
}

// NewKafkaConfig Конструктор конфига для подключения к Kafka.
func NewKafkaConfig(os OS) (*KafkaConfig, error) {
	if os == nil {
		return nil, fmt.Errorf("os must not be nil")
	}

	if firstHost := os.Getenv(kafkaHostEnvName + "1"); firstHost == "" {
		return &KafkaConfig{
			addresses: []string{},
			enabled:   false,
		}, nil
	}

	addresses := make([]string, 0)
	i := 1
	for {
		hostPort := os.Getenv(kafkaHostEnvName + strconv.Itoa(i))
		if hostPort == "" {
			break
		}

		addresses = append(addresses, hostPort)

		i++
	}

	return &KafkaConfig{
		addresses: addresses,
		enabled:   true,
	}, nil
}

// Addresses ...
func (cfg *KafkaConfig) Addresses() []string {
	return cfg.addresses
}

// Enabled ...
func (cfg *KafkaConfig) Enabled() bool {
	return cfg.enabled
}
