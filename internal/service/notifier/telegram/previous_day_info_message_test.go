package telegram

import (
	"errors"
	"fmt"
	"testing"
	"time"

	TGBotAPI "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var messagePreviousDayInfo = TGBotAPI.MessageConfig{
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
	Text:                  "Вчера Вы были в офисе с 08:10 до 17:05\nУчтенное время 8 ч. 10 м.",
	ParseMode:             "",
	DisableWebPagePreview: false,
}

func TestSendPreviousDayInfoMessage(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)

	tests := map[string]struct {
		bot              func() botAPI
		metrics          func() metrics
		apiPBClient      func() apiPBClient
		aggregatorClient func() aggregatorClient
		params           ParamsPreviousDayInfo

		expectedNilObj bool
		expectedErr    error
	}{
		"send message": {
			bot: func() botAPI {
				botMock := NewBotAPIMock(mc)
				require.NotNil(t, botMock)

				botMock.SendMock.Expect(messagePreviousDayInfo).Times(1).Return(TGBotAPI.Message{}, nil)

				return botMock
			},
			metrics: func() metrics {
				metricsMock := NewMetricsMock(mc)
				require.NotNil(t, metricsMock)

				metricsMock.IncMessageCounterMock.Expect().Times(1)

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
			params: ParamsPreviousDayInfo{
				TelegramID:   43403446,
				SecondsStart: time.Date(2025, time.May, 10, 8, 10, 43, 0, time.Local),
				SecondsEnd:   time.Date(2025, time.May, 10, 17, 5, 30, 0, time.Local),
				Hours:        8,
				Minutes:      10,
				Breaks:       "",
			},

			expectedNilObj: false,
			expectedErr:    nil,
		},
		"send error message": {
			bot: func() botAPI {
				botMock := NewBotAPIMock(mc)
				require.NotNil(t, botMock)

				botMock.
					SendMock.
					Expect(messagePreviousDayInfo).
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
			params: ParamsPreviousDayInfo{
				TelegramID:   43403446,
				SecondsStart: time.Date(2025, time.May, 10, 8, 10, 43, 0, time.Local),
				SecondsEnd:   time.Date(2025, time.May, 10, 17, 5, 30, 0, time.Local),
				Hours:        8,
				Minutes:      10,
				Breaks:       "",
			},

			expectedNilObj: false,
			expectedErr:    fmt.Errorf("error send message: %w", errors.New("API error")),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			n, err := NewTelegramNotifier(test.bot(), test.metrics(), test.apiPBClient(), test.aggregatorClient())
			assert.NotNil(t, n)
			assert.Nil(t, err)

			err = n.SendPreviousDayInfoMessage(minimock.AnyContext, test.params)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
