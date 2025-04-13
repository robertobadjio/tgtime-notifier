package telegram

import (
	"errors"
	"fmt"
	"testing"

	TGBotAPI "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/robertobadjio/tgtime-notifier/internal/metric"
	notifier2 "github.com/robertobadjio/tgtime-notifier/internal/service/notifier"
	command2 "github.com/robertobadjio/tgtime-notifier/internal/service/notifier/telegram/command"
)

var update = TGBotAPI.Update{
	UpdateID: 1,
	Message: &TGBotAPI.Message{
		From: &TGBotAPI.User{
			ID: 43403446,
		},
	},
}

var messageUpdateHandler = TGBotAPI.MessageConfig{
	BaseChat: TGBotAPI.BaseChat{
		ChatID: 43403446,
		ReplyMarkup: TGBotAPI.ReplyKeyboardMarkup{
			Keyboard: [][]TGBotAPI.KeyboardButton{{
				TGBotAPI.KeyboardButton{Text: "⏳ Рабочее время"},
				TGBotAPI.KeyboardButton{Text: "🗓 Статистика за рабочий период"},
			}},
			ResizeKeyboard: true,
		},
	},
	Text:                  "Тестовое сообщение",
	ParseMode:             "",
	DisableWebPagePreview: false,
}

func Test(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)

	tests := map[string]struct {
		bot     func() BotAPIInterface
		factory func() command2.Factory
		metrics func() metric.Metrics
		params  notifier2.Params

		expectedNilObj bool
		expectedErr    error
	}{
		"send error message, error cast interface param": {
			bot: func() BotAPIInterface {
				botAPIMock := NewBotAPIInterfaceMock(mc)
				require.NotNil(t, botAPIMock)

				return botAPIMock
			},
			factory: func() command2.Factory {
				factoryMock := command2.NewFactoryMock(mc)
				require.NotNil(t, factoryMock)

				return factoryMock
			},
			metrics: func() metric.Metrics {
				metricsMock := metric.NewMetricsMock(mc)
				require.NotNil(t, metricsMock)

				return metricsMock
			},
			params: ParamsPreviousDayInfo{TelegramID: 43403446},

			expectedErr: errors.New("error cast interface param"),
		},
		"send message": {
			bot: func() BotAPIInterface {
				botMock := NewBotAPIInterfaceMock(mc)
				require.NotNil(t, botMock)

				botMock.
					SendMock.
					Expect(messageUpdateHandler).
					Times(1).
					Return(TGBotAPI.Message{}, nil)

				return botMock
			},
			factory: func() command2.Factory {
				factoryMock := command2.NewFactoryMock(mc)
				require.NotNil(t, factoryMock)

				messageMock := command2.NewCommandMock(mc)
				require.NotNil(t, messageMock)

				messageMock.
					GetMessageMock.
					Expect(minimock.AnyContext).
					Times(1).
					Return("Тестовое сообщение", nil)

				factoryMock.GetCommandHandlerMock.Expect(update).Times(1).Return(messageMock)

				return factoryMock
			},
			metrics: func() metric.Metrics {
				metricMock := metric.NewMetricsMock(mc)
				require.NotNil(t, metricMock)

				metricMock.IncMessageCounterMock.Expect().Times(1)

				return metricMock
			},
			params: ParamsUpdate{Update: update},

			expectedNilObj: false,
			expectedErr:    nil,
		},
		"send message fail": {
			bot: func() BotAPIInterface {
				botMock := NewBotAPIInterfaceMock(mc)
				require.NotNil(t, botMock)

				botMock.
					SendMock.
					Expect(messageUpdateHandler).
					Times(1).
					Return(TGBotAPI.Message{}, errors.New("API error"))

				return botMock
			},
			factory: func() command2.Factory {
				factoryMock := command2.NewFactoryMock(mc)
				require.NotNil(t, factoryMock)

				messageMock := command2.NewCommandMock(mc)
				require.NotNil(t, messageMock)

				messageMock.
					GetMessageMock.
					Expect(minimock.AnyContext).
					Times(1).
					Return("Тестовое сообщение", nil)

				factoryMock.GetCommandHandlerMock.Expect(update).Times(1).Return(messageMock)

				return factoryMock
			},
			metrics: func() metric.Metrics {
				metricsMock := metric.NewMetricsMock(mc)
				require.NotNil(t, metricsMock)

				return metricsMock
			},
			params: ParamsUpdate{Update: update},

			expectedNilObj: false,
			expectedErr:    fmt.Errorf("error send message: %w", errors.New("API error")),
		},
		"send message with internal error": {
			bot: func() BotAPIInterface {
				botMock := NewBotAPIInterfaceMock(mc)
				require.NotNil(t, botMock)

				return botMock
			},
			factory: func() command2.Factory {
				factoryMock := command2.NewFactoryMock(mc)
				require.NotNil(t, factoryMock)

				messageMock := command2.NewCommandMock(mc)
				require.NotNil(t, messageMock)

				messageMock.
					GetMessageMock.
					Expect(minimock.AnyContext).
					Times(1).
					Return("", errors.New("error getting user by telegram ID"))

				factoryMock.GetCommandHandlerMock.Expect(update).Times(1).Return(messageMock)

				return factoryMock
			},
			metrics: func() metric.Metrics {
				metricsMock := metric.NewMetricsMock(mc)
				require.NotNil(t, metricsMock)

				return metricsMock
			},
			params: ParamsUpdate{Update: update},

			expectedNilObj: false,
			expectedErr:    fmt.Errorf("error getting text message: %w", errors.New("error getting user by telegram ID")),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			n, err := NewTelegramNotifier(test.bot(), test.factory(), test.metrics())
			assert.NotNil(t, n)
			assert.Nil(t, err)

			err = n.SendCommandMessage(minimock.AnyContext, test.params)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
