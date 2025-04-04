package config

import (
	"fmt"
	"strconv"
)

const (
	kafkaHostEnvName = "KAFKA_BROKER_"
)

// KafkaConfig Конфиг для подключения к Kafka
type KafkaConfig interface {
	GetAddresses() []string
}

type kafkaConfig struct {
	addresses []string
}

// NewKafkaConfig Конструктор конфига для подключения к Kafka
func NewKafkaConfig(os OS) (KafkaConfig, error) {
	if os == nil {
		return nil, fmt.Errorf("os must not be nil")
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

	return &kafkaConfig{
		addresses: addresses,
	}, nil
}

// GetAddresses ...
func (cfg *kafkaConfig) GetAddresses() []string {
	return cfg.addresses
}
