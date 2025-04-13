package telegram

import (
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"github.com/robertobadjio/tgtime-notifier/internal/metric"
	command2 "github.com/robertobadjio/tgtime-notifier/internal/service/notifier/telegram/command"
)

func TestTelegramNotifier_New(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)

	tests := map[string]struct {
		bot     func() BotAPIInterface
		factory func() command2.Factory
		metrics func() metric.Metrics

		expectedNilObj bool
		expectedErr    error
	}{
		"telegram bot is nil": {
			bot:            func() BotAPIInterface { return nil },
			factory:        func() command2.Factory { return nil },
			metrics:        func() metric.Metrics { return nil },
			expectedNilObj: true,
			expectedErr:    errors.New("telegram bot is nil"),
		},
		"telegram factory is nil": {
			bot: func() BotAPIInterface {
				return NewBotAPIInterfaceMock(mc)
			},
			factory:        func() command2.Factory { return nil },
			metrics:        func() metric.Metrics { return nil },
			expectedNilObj: true,
			expectedErr:    errors.New("telegram factory is nil"),
		},
		"metrics is nil": {
			bot: func() BotAPIInterface {
				return NewBotAPIInterfaceMock(mc)
			},
			factory: func() command2.Factory {
				return command2.NewFactoryMock(mc)
			},
			metrics: func() metric.Metrics {
				return nil
			},
			expectedNilObj: true,
			expectedErr:    errors.New("metrics is nil"),
		},
		"create telegram notifier": {
			bot: func() BotAPIInterface {
				return NewBotAPIInterfaceMock(mc)
			},
			factory: func() command2.Factory {
				return command2.NewFactoryMock(mc)
			},
			metrics: func() metric.Metrics {
				return metric.NewMetricsMock(mc)
			},
			expectedNilObj: false,
			expectedErr:    nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			n, err := NewTelegramNotifier(test.bot(), test.factory(), test.metrics())
			assert.Equal(t, test.expectedErr, err)

			if test.expectedNilObj {
				assert.Nil(t, n)
			} else {
				assert.NotNil(t, n)
			}
		})
	}
}
