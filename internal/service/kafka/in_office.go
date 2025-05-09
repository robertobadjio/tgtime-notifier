package kafka

import (
	"context"
	"errors"
	"fmt"
	"io"

	kafkaLib "github.com/segmentio/kafka-go"

	"github.com/robertobadjio/tgtime-notifier/internal/logger"
	"github.com/robertobadjio/tgtime-notifier/internal/service/notifier/telegram"
)

const inOfficeTopic = "in-office"

const partition = 0

// ConsumeInOffice Чтение сообщений из кафки о приходе сотрудника в офис / на работу.
func (k *Kafka) ConsumeInOffice(ctx context.Context) error {
	r := k.buildReader(inOfficeTopic)
	defer func() {
		if err := r.Close(); err != nil {
			logger.Warn(
				"component", "kafka",
				"during", "consume in office",
				"desc", "failed to close reader",
				"error", err.Error(),
			)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("exit from consumer in office: %w", ctx.Err())
		default:
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("exit from consumer in office: %w", ctx.Err())
		default:
			m, err := r.ReadMessage(ctx)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break // TODO: ?
				}
				return fmt.Errorf("reading message: %w", err)
			}
			userResponse, err := k.tgTimeAPIClient.GetUserByMacAddress(ctx, string(m.Value))
			if err != nil {
				logger.Error(
					"component", "kafka",
					"during", "read",
					"topic", inOfficeTopic,
					"desc", "error getting user by mac address "+string(m.Value),
					"error", err.Error(),
				)
				continue
			}

			err = k.notifier.SendWelcomeMessage(
				ctx,
				telegram.ParamsWelcomeMessage{TelegramID: userResponse.User.TelegramId},
			)
			if err != nil {
				logger.Error(
					"component", "kafka",
					"during", "read",
					"topic", inOfficeTopic,
					"desc", "error sending welcome message",
					"error", err.Error(),
				)
			}
		}
	}
}

func (k *Kafka) buildReader(topicName string) *kafkaLib.Reader {
	return kafkaLib.NewReader(kafkaLib.ReaderConfig{
		Brokers:   k.brokers,
		Topic:     topicName,
		Partition: partition,
		GroupID:   "",
		MaxBytes:  10e3,
	})
}
