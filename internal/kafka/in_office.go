package kafka

import (
	"context"
	"errors"
	"fmt"
	"io"

	kafkaLib "github.com/segmentio/kafka-go"

	"github.com/robertobadjio/tgtime-notifier/internal/logger"
	"github.com/robertobadjio/tgtime-notifier/internal/notifier/telegram"
)

const inOfficeTopic = "in-office"

const partition = 0

// ConsumeInOffice Чтение сообщений из кафки о приходе сотрудника в офис / на работу
func (k *Kafka) ConsumeInOffice(ctx context.Context) error {
	r := k.buildReader(inOfficeTopic)
	defer func() {
		if err := r.Close(); err != nil {
			logger.Log("failed to close reader:", err)
		}
	}()

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return fmt.Errorf("reading message: %w", err)
		}
		userResponse, err := k.tgTimeAPIClient.GetUserByMacAddress(ctx, string(m.Value))
		if err != nil {
			logger.Log("kafka", "read", "topic", inOfficeTopic, "error", err.Error(), "desc", "error getting user by mac address "+string(m.Value))
			continue
		}

		err = k.notifier.SendWelcomeMessage(
			ctx,
			telegram.ParamsWelcomeMessage{TelegramID: userResponse.User.TelegramId},
		)
		if err != nil {
			logger.Log("kafka", "read", "topic", inOfficeTopic, "error", err.Error())
		}
	}

	return nil
}

func (k *Kafka) buildReader(topicName string) *kafkaLib.Reader {
	return kafkaLib.NewReader(kafkaLib.ReaderConfig{
		Brokers:   k.config.GetAddresses(),
		Topic:     topicName,
		Partition: partition,
		GroupID:   "",
		MaxBytes:  10e3,
	})
}
