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

	"github.com/robertobadjio/tgtime-notifier/internal/metric"
	notifier2 "github.com/robertobadjio/tgtime-notifier/internal/service/notifier"
	command2 "github.com/robertobadjio/tgtime-notifier/internal/service/notifier/telegram/command"
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
		bot     func() BotAPIInterface
		factory func() command2.Factory
		metrics func() metric.Metrics
		params  notifier2.Params

		expectedNilObj bool
		expectedErr    error
	}{
		"send error message, error cast interface param": {
			bot: func() BotAPIInterface {
				botMock := NewBotAPIInterfaceMock(mc)
				require.NotNil(t, botMock)

				return botMock
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
			params: ParamsWelcomeMessage{TelegramID: 43403446},

			expectedErr: errors.New("error cast interface param"),
		},
		"send message": {
			bot: func() BotAPIInterface {
				botMock := NewBotAPIInterfaceMock(mc)
				require.NotNil(t, botMock)

				botMock.SendMock.Expect(messagePreviousDayInfo).Times(1).Return(TGBotAPI.Message{}, nil)

				return botMock
			},
			factory: func() command2.Factory {
				factoryMock := command2.NewFactoryMock(mc)
				require.NotNil(t, factoryMock)

				return factoryMock
			},
			metrics: func() metric.Metrics {
				metricsMock := metric.NewMetricsMock(mc)
				require.NotNil(t, metricsMock)

				metricsMock.IncMessageCounterMock.Expect().Times(1)

				return metricsMock
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
			bot: func() BotAPIInterface {
				botMock := NewBotAPIInterfaceMock(mc)
				require.NotNil(t, botMock)

				botMock.
					SendMock.
					Expect(messagePreviousDayInfo).
					Times(1).
					Return(TGBotAPI.Message{}, errors.New("API error"))

				return botMock
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

			n, err := NewTelegramNotifier(test.bot(), test.factory(), test.metrics())
			assert.NotNil(t, n)
			assert.Nil(t, err)

			err = n.SendPreviousDayInfoMessage(minimock.AnyContext, test.params)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
