package kafka

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/go-kit/kit/log"
	"github.com/robertobadjio/tgtime-notifier/internal/api_pb"
	"github.com/robertobadjio/tgtime-notifier/internal/config"
	"github.com/robertobadjio/tgtime-notifier/internal/notifier/telegram"
	kafkaLib "github.com/segmentio/kafka-go"
)

// Kafka Клиент для подключения к кафке
type Kafka struct {
	logger    log.Logger
	addresses []string
}

// NewKafka Конструктор клиента
func NewKafka(logger log.Logger, addresses []string) *Kafka {
	return &Kafka{logger: logger, addresses: addresses}
}

// ConsumeInOffice Чтение сообщений из кафки о приходе сотрудника в офис / на работу
func (k *Kafka) ConsumeInOffice(ctx context.Context, tn *telegram.Notifier) error {
	r := k.buildReader(inOfficeTopic)
	defer func() {
		if err := r.Close(); err != nil {
			_ = k.logger.Log("failed to close reader:", err)
		}
	}()

	tgTimeAPIConfig, _ := config.NewTgTimeAPIConfig() // TODO: ?!

	tgtimeClient := api_pb.NewClient(tgTimeAPIConfig, k.logger)

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return fmt.Errorf("reading message: %w", err)
		}
		//fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
		//fmt.Printf(string(m.Value))

		userResponse, err := tgtimeClient.GetUserByMacAddress(ctx, string(m.Value))
		if err != nil {
			fmt.Println("error getting user by mac address " + string(m.Value))
			continue
		}

		fmt.Printf("%+v\n", userResponse.User.TelegramId)
		err = tn.SendWelcomeMessage(ctx, userResponse.User.TelegramId)
		if err != nil {
			fmt.Println("error sending welcome message: ", err.Error())
		}
	}

	return nil
}

func (k *Kafka) buildReader(topicName string) *kafkaLib.Reader {
	return kafkaLib.NewReader(kafkaLib.ReaderConfig{
		Brokers:   k.addresses,
		Topic:     topicName,
		Partition: partition,
		GroupID:   "",
		MaxBytes:  10e3,
	})
}
