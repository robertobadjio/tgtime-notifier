package command

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/robertobadjio/tgtime-notifier/internal/service/client/aggregator"
	"github.com/robertobadjio/tgtime-notifier/internal/service/client/api_pb"
)

// Factory ...
type Factory interface {
	GetCommandHandler(update tgbotapi.Update) Command
}

type factory struct {
	TGTimeAPIClient        api_pb.Client
	TGTimeAggregatorClient aggregator.Client
}

// NewFactory ...
func NewFactory(TGTimeAPIClient api_pb.Client, TGTimeAggregatorClient aggregator.Client) Factory {
	return &factory{
		TGTimeAPIClient:        TGTimeAPIClient,
		TGTimeAggregatorClient: TGTimeAggregatorClient,
	}
}

// GetCommandHandler ...
func (f *factory) GetCommandHandler(update tgbotapi.Update) Command {
	switch Type(update.Message.Text) {
	case ButtonStart:
		return NewStartCommand()
	case ButtonWorkingTime:
		return NewWorkingTimeCommand(f.TGTimeAPIClient, f.TGTimeAggregatorClient, int64(update.Message.From.ID))
	case ButtonStatCurrentWorkingPeriod:
		return NewStatCurrentWorkingPeriodCommand()
	default:
		return NewUnknownCommand()
	}
}
