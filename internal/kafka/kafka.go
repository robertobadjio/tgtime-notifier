package kafka

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kit/kit/log"
	kafkaLib "github.com/segmentio/kafka-go"
	"io"
	"tgtime-notifier/internal/notifier"
)

type Kafka struct {
	logger log.Logger
	host   string
	port   string
}

func NewKafka(logger log.Logger, host, port string) *Kafka {
	return &Kafka{logger: logger, host: host, port: port}
}

func (k *Kafka) buildReader(topicName string) *kafkaLib.Reader {
	return kafkaLib.NewReader(kafkaLib.ReaderConfig{
		Brokers:   []string{buildAddress(k.host, k.port)},
		Topic:     topicName,
		Partition: partition,
		GroupID:   "",
		MaxBytes:  10e3,
	})
}

func (k *Kafka) ConsumeInOffice(ctx context.Context, nt notifier.Notifier) error {
	r := k.buildReader(inOfficeTopic)
	defer func() {
		if err := r.Close(); err != nil {
			_ = k.logger.Log("failed to close reader:", err)
		}
	}()

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				return fmt.Errorf("reading message: %w", err)
			}
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))

		// TODO: Вынести отправку сообщения
		//_ = nt.SendWelcomeMessage(ctx, int64(telegramId)) // 343536263 // TODO: Handle error
		err = nt.SendWelcomeMessage(ctx, 343536263)
		if err != nil {
			return fmt.Errorf("sending welcome message: %w", err)
		}
	}

	return nil
}

func (k *Kafka) ConsumePreviousDayInfo(ctx context.Context, nt notifier.Notifier) error {
	// TODO: !
	panic("implement me")
}

func buildAddress(host, port string) string {
	return host + ":" + port
}
