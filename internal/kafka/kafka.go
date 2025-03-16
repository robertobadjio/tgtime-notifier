package kafka

import (
	"github.com/robertobadjio/tgtime-notifier/internal/config"
	"github.com/robertobadjio/tgtime-notifier/internal/service/client/api_pb"
	"github.com/robertobadjio/tgtime-notifier/internal/service/notifier"
)

// Kafka Клиент для подключения к кафке
type Kafka struct {
	notifier        notifier.Notifier
	config          config.KafkaConfig
	tgTimeAPIClient api_pb.Client
}

// NewKafka Конструктор клиента
func NewKafka(
	config config.KafkaConfig,
	notifier notifier.Notifier,
	tgTimeAPIClient api_pb.Client,
) *Kafka {
	return &Kafka{config: config, notifier: notifier, tgTimeAPIClient: tgTimeAPIClient}
}
