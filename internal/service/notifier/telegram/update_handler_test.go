package telegram

import (
	"errors"
	"fmt"
	"testing"

	TGBotAPI "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var updateStart = TGBotAPI.Update{
	UpdateID: 1,
	Message: &TGBotAPI.Message{
		From: &TGBotAPI.User{
			ID: 43403446,
		},
		Text: ButtonStart,
	},
}

var updateWorkingTime = TGBotAPI.Update{
	UpdateID: 1,
	Message: &TGBotAPI.Message{
		From: &TGBotAPI.User{
			ID: 43403446,
		},
		Text: ButtonWorkingTime,
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
	Text:                  "Добро пожаловать. Используйте кнопки для получения информации.",
	ParseMode:             "",
	DisableWebPagePreview: false,
}

func Test_SendCommandMessage(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)

	tests := map[string]struct {
		bot              func() botAPI
		metrics          func() metrics
		apiPBClient      func() apiPBClient
		aggregatorClient func() aggregatorClient
		params           ParamsUpdate

		expectedNilObj bool
		expectedErr    error
	}{
		"send message": {
			bot: func() botAPI {
				botMock := NewBotAPIMock(mc)
				require.NotNil(t, botMock)

				botMock.
					SendMock.
					Expect(messageUpdateHandler).
					Times(1).
					Return(TGBotAPI.Message{}, nil)

				return botMock
			},
			metrics: func() metrics {
				metricMock := NewMetricsMock(mc)
				require.NotNil(t, metricMock)

				metricMock.IncMessageCounterMock.Expect().Times(1)

				return metricMock
			},
			apiPBClient: func() apiPBClient {
				clientMock := NewApiPBClientMock(mc)
				require.NotNil(t, clientMock)

				return clientMock
			},
			aggregatorClient: func() aggregatorClient {
				clientMock := NewAggregatorClientMock(mc)
				require.NotNil(t, clientMock)

				return clientMock
			},
			params: ParamsUpdate{Update: updateStart},

			expectedNilObj: false,
			expectedErr:    nil,
		},
		"send message fail: telegram bot API error": {
			bot: func() botAPI {
				botMock := NewBotAPIMock(mc)
				require.NotNil(t, botMock)

				botMock.
					SendMock.
					Expect(messageUpdateHandler).
					Times(1).
					Return(TGBotAPI.Message{}, errors.New("API error"))

				return botMock
			},
			metrics: func() metrics {
				metricsMock := NewMetricsMock(mc)
				require.NotNil(t, metricsMock)

				return metricsMock
			},
			apiPBClient: func() apiPBClient {
				clientMock := NewApiPBClientMock(mc)
				require.NotNil(t, clientMock)

				return clientMock
			},
			aggregatorClient: func() aggregatorClient {
				clientMock := NewAggregatorClientMock(mc)
				require.NotNil(t, clientMock)

				return clientMock
			},
			params: ParamsUpdate{Update: updateStart},

			expectedNilObj: false,
			expectedErr:    fmt.Errorf("error send message: %w", errors.New("API error")),
		},
		"send message with internal error": {
			bot: func() botAPI {
				botMock := NewBotAPIMock(mc)
				require.NotNil(t, botMock)

				return botMock
			},
			metrics: func() metrics {
				metricsMock := NewMetricsMock(mc)
				require.NotNil(t, metricsMock)

				return metricsMock
			},
			apiPBClient: func() apiPBClient {
				clientMock := NewApiPBClientMock(mc)
				require.NotNil(t, clientMock)

				clientMock.
					GetUserByTelegramIDMock.
					Expect(minimock.AnyContext, int64(43403446)).
					Times(1).
					Return(nil, errors.New("internal error"))

				return clientMock
			},
			aggregatorClient: func() aggregatorClient {
				clientMock := NewAggregatorClientMock(mc)
				require.NotNil(t, clientMock)

				return clientMock
			},
			params: ParamsUpdate{Update: updateWorkingTime},

			expectedNilObj: false,
			expectedErr: fmt.Errorf(
				"handle working time command: %w",
				fmt.Errorf("error getting user by telegram ID: %w", errors.New("internal error")),
			),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			n, err := NewTelegramNotifier(test.bot(), test.metrics(), test.apiPBClient(), test.aggregatorClient())
			assert.NotNil(t, n)
			assert.Nil(t, err)

			err = n.SendCommandMessage(minimock.AnyContext, test.params)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
