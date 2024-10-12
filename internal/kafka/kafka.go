package kafka

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kit/kit/log"
	kafkaLib "github.com/segmentio/kafka-go"
	"io"
	"tgtime-notifier/internal/api_pb"
	"tgtime-notifier/internal/config"
	"tgtime-notifier/internal/notifier/telegram"
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

func (k *Kafka) ConsumeInOffice(ctx context.Context, tn *telegram.TelegramNotifier) error {
	r := k.buildReader(inOfficeTopic)
	defer func() {
		if err := r.Close(); err != nil {
			_ = k.logger.Log("failed to close reader:", err)
		}
	}()

	tgtimeClient := api_pb.NewClient(*config.New(), k.logger)

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				return fmt.Errorf("reading message: %w", err)
			}
		}
		//fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
		//fmt.Printf(string(m.Value))

		userResponse, err := tgtimeClient.GetUserByMacAddress(ctx, string(m.Value))
		if err != nil {
			return fmt.Errorf("error getting user by mac address: %w", err)
		}
		if userResponse.User == nil {
			fmt.Println("user not found with mac address " + string(m.Value))
			continue
		}

		err = tn.SendWelcomeMessage(ctx, userResponse.User.TelegramId)
		if err != nil {
			fmt.Println("error sending welcome message: ", err.Error())
		}
	}

	return nil
}

func buildAddress(host, port string) string {
	return host + ":" + port
}
