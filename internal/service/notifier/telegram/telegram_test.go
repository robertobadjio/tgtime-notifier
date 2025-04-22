package telegram

import (
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTelegramNotifier_New(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)

	tests := map[string]struct {
		bot              func() botAPI
		metrics          func() metrics
		apiPBClient      func() apiPBClient
		aggregatorClient func() aggregatorClient

		expectedNilObj bool
		expectedErr    error
	}{
		"telegram bot is nil": {
			bot: func() botAPI {
				return nil
			},
			metrics: func() metrics {
				return nil
			},
			apiPBClient: func() apiPBClient {
				return nil
			},
			aggregatorClient: func() aggregatorClient {
				return nil
			},

			expectedNilObj: true,
			expectedErr:    errors.New("telegram bot must be set"),
		},
		"metrics is nil": {
			bot: func() botAPI {
				botMock := NewBotAPIMock(mc)
				require.NotNil(t, botMock)

				return botMock
			},
			metrics: func() metrics {
				return nil
			},
			apiPBClient: func() apiPBClient {
				return nil
			},
			aggregatorClient: func() aggregatorClient {
				return nil
			},

			expectedNilObj: true,
			expectedErr:    errors.New("metrics must be set"),
		},
		"api pb client is nil": {
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
				return nil
			},
			aggregatorClient: func() aggregatorClient {
				return nil
			},

			expectedNilObj: true,
			expectedErr:    errors.New("TGTimeAPIClient must be set"),
		},
		"aggregator client is nil": {
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

				return clientMock
			},
			aggregatorClient: func() aggregatorClient {
				return nil
			},

			expectedNilObj: true,
			expectedErr:    errors.New("TGTimeAggregatorClient must be set"),
		},
		"create telegram TGNotifier": {
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

				return clientMock
			},
			aggregatorClient: func() aggregatorClient {
				clientMock := NewAggregatorClientMock(mc)
				require.NotNil(t, clientMock)

				return clientMock
			},

			expectedNilObj: false,
			expectedErr:    nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			n, err := NewTelegramNotifier(test.bot(), test.metrics(), test.apiPBClient(), test.aggregatorClient())
			assert.Equal(t, test.expectedErr, err)

			if test.expectedNilObj {
				assert.Nil(t, n)
			} else {
				assert.NotNil(t, n)
			}
		})
	}
}
