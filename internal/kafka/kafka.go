package kafka

import (
	"context"

	"github.com/robertobadjio/tgtime-notifier/internal/config"
	"github.com/robertobadjio/tgtime-notifier/internal/service/client/api_pb"
	"github.com/robertobadjio/tgtime-notifier/internal/service/notifier/telegram"
)

type notifier interface {
	SendWelcomeMessage(ctx context.Context, params telegram.ParamsWelcomeMessage) error
}

// Kafka Клиент для подключения к кафке.
type Kafka struct {
	notifier        notifier
	config          config.KafkaConfig
	tgTimeAPIClient api_pb.Client
}

// NewKafka Конструктор клиента.
func NewKafka(
	config config.KafkaConfig,
	notifier notifier,
	tgTimeAPIClient api_pb.Client,
) *Kafka {
	return &Kafka{config: config, notifier: notifier, tgTimeAPIClient: tgTimeAPIClient}
}
