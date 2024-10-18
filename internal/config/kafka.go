package config

import (
	"os"
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
func NewKafkaConfig() (KafkaConfig, error) {
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

func (cfg *kafkaConfig) GetAddresses() []string {
	return cfg.addresses
}
