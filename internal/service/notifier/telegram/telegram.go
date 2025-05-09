package telegram

import (
	"context"
	"errors"
	"fmt"

	TGBotAPI "github.com/go-telegram-bot-api/telegram-bot-api"
	pb "github.com/robertobadjio/tgtime-aggregator/pkg/api/time_v1"
	pbapiv1 "github.com/robertobadjio/tgtime-api/api/v1/pb/api"
)

// botAPI ...
type botAPI interface {
	Send(c TGBotAPI.Chattable) (TGBotAPI.Message, error)
}

type metrics interface {
	IncMessageCounter()
}

type aggregatorClient interface {
	GetTimeSummary(
		ctx context.Context,
		macAddress, date string,
	) (*pb.GetSummaryResponse, error)
}

type apiPBClient interface {
	GetUserByTelegramID(
		ctx context.Context,
		telegramID int64,
	) (*pbapiv1.GetUserByTelegramIdResponse, error)
}

// TGNotifier Telegram-нотификатор.
type TGNotifier struct {
	bot                    botAPI
	metrics                metrics
	TGTimeAPIClient        apiPBClient
	TGTimeAggregatorClient aggregatorClient
}

// NewTelegramNotifier Конструктор для создания Telegram-нотификатора.
func NewTelegramNotifier(
	bot botAPI,
	metrics metrics,
	TGTimeAPIClient apiPBClient,
	TGTimeAggregatorClient aggregatorClient,
) (*TGNotifier, error) {
	if bot == nil {
		return nil, errors.New("telegram bot must be set")
	}

	if metrics == nil {
		return nil, errors.New("metrics must be set")
	}

	if TGTimeAPIClient == nil {
		return nil, errors.New("TGTimeAPIClient must be set")
	}

	if TGTimeAggregatorClient == nil {
		return nil, errors.New("TGTimeAggregatorClient must be set")
	}

	return &TGNotifier{
		bot:                    bot,
		metrics:                metrics,
		TGTimeAPIClient:        TGTimeAPIClient,
		TGTimeAggregatorClient: TGTimeAggregatorClient,
	}, nil
}

func (tn *TGNotifier) setKeyboard(message TGBotAPI.MessageConfig) TGBotAPI.MessageConfig {
	message.ReplyMarkup = TGBotAPI.NewReplyKeyboard(
		TGBotAPI.NewKeyboardButtonRow(
			TGBotAPI.NewKeyboardButton(ButtonWorkingTime),
			TGBotAPI.NewKeyboardButton(ButtonStatCurrentWorkingPeriod),
		),
	)

	return message
}

func (tn *TGNotifier) sendMessage(text string, telegramID int64) error {
	var err error
	if telegramID > 0 {
		_, err = tn.bot.Send(tn.setKeyboard(TGBotAPI.NewMessage(
			telegramID,
			text,
		)))
	} else {
		_, err = tn.bot.Send(TGBotAPI.NewMessage(
			telegramID,
			text,
		))
	}

	if err != nil {
		return fmt.Errorf("error send message: %w", err)
	}

	tn.metrics.IncMessageCounter()

	return nil
}
