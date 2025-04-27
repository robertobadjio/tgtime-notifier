package kafka

import (
	"context"

	"github.com/robertobadjio/tgtime-notifier/internal/service/client/api_pb"
	"github.com/robertobadjio/tgtime-notifier/internal/service/notifier/telegram"
)

type notifier interface {
	SendWelcomeMessage(ctx context.Context, params telegram.ParamsWelcomeMessage) error
}

// Kafka Клиент для подключения к кафке.
type Kafka struct {
	notifier        notifier
	brokers         []string
	tgTimeAPIClient *api_pb.Client
}

// NewKafka Конструктор клиента.
func NewKafka(
	brokers []string,
	notifier notifier,
	tgTimeAPIClient *api_pb.Client,
) *Kafka {
	return &Kafka{brokers: brokers, notifier: notifier, tgTimeAPIClient: tgTimeAPIClient}
}
