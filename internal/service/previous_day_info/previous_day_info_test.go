package previous_day_info

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pb "github.com/robertobadjio/tgtime-aggregator/pkg/api/time_v1"
	pbapiv1 "github.com/robertobadjio/tgtime-api/api/v1/pb/api"
	"github.com/robertobadjio/tgtime-notifier/internal/helper"
	"github.com/robertobadjio/tgtime-notifier/internal/service/notifier/telegram"
)

func TestCancelContext(t *testing.T) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("context cancelled", func(t *testing.T) {
		innerCtx, innerCancel := context.WithCancel(ctx)

		controller := minimock.NewController(t)

		notifierMock := NewNotifierMock(controller)
		require.NotNil(t, notifierMock)

		aggregatorClientMock := NewAggregatorClientMock(controller)
		require.NotNil(t, aggregatorClientMock)

		apiClientMock := NewAPIClientMock(controller)
		require.NotNil(t, apiClientMock)

		previousDayInfoService, err := NewPreviousDayInfo(
			apiClientMock,
			aggregatorClientMock,
			notifierMock,
			clockwork.NewFakeClock(),
			12,
			0,
			0,
		)
		require.Nil(t, err)
		require.NotNil(t, previousDayInfoService)

		handler := func(_ context.Context, _ string) error {
			return nil
		}

		innerCancel()
		err = previousDayInfoService.everyDayByHour(innerCtx, handler, time.NewTicker(time.Minute))
		assert.ErrorIs(t, err, context.Canceled)
	})
}

func TestTick(t *testing.T) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("handle function in defined time", func(t *testing.T) {
		innerCtx, innerCancel := context.WithCancel(ctx)
		controller := minimock.NewController(t)

		notifierMock := NewNotifierMock(controller)
		require.NotNil(t, notifierMock)

		aggregatorClientMock := NewAggregatorClientMock(controller)
		require.NotNil(t, aggregatorClientMock)

		apiClientMock := NewAPIClientMock(controller)
		require.NotNil(t, apiClientMock)

		clockMock := clockwork.NewFakeClock()
		timeHandle := time.Now().Add(time.Second)

		previousDayInfoService, err := NewPreviousDayInfo(
			apiClientMock,
			aggregatorClientMock,
			notifierMock,
			clockMock,
			timeHandle.Hour(),
			timeHandle.Minute(),
			timeHandle.Second(),
		)
		require.Nil(t, err)
		require.NotNil(t, previousDayInfoService)

		const resultText string = "handler called"

		handler := func(_ context.Context, _ string) error {
			fmt.Print(resultText)
			return nil
		}

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			output := captureOutput(func() {
				_ = previousDayInfoService.everyDayByHour(innerCtx, handler, time.NewTicker(time.Second))
			})
			assert.Equal(t, resultText, output)
			wg.Done()
		}()

		clockMock.Advance(time.Second)
		time.Sleep(time.Second)

		innerCancel()

		wg.Wait()
	})
}

func captureOutput(f func()) string {
	orig := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	os.Stdout = orig
	_ = w.Close()
	out, _ := io.ReadAll(r)
	return string(out)
}

var expectedTimeSummaryByUsers = []*pb.Summary{
	0: {
		MacAddress:   "00:1b:63:84:45:e6",
		Seconds:      123123,
		Breaks:       nil,
		Date:         "2025-05-06",
		SecondsStart: 1746548943,
		SecondsEnd:   1746548950,
	},
}

func Test_sendAllUsersNotify(t *testing.T) {
	t.Parallel()

	controller := minimock.NewController(t)

	tests := map[string]struct {
		notifier         func() notifier
		aggregatorClient func() aggregatorClient
		APIClient        func() APIClient
		date             string

		expectedErr error
	}{
		"error get time summary for all users": {
			notifier: func() notifier {
				notifierMock := NewNotifierMock(controller)
				require.NotNil(t, notifierMock)

				return notifierMock
			},
			aggregatorClient: func() aggregatorClient {
				aggregatorClientMock := NewAggregatorClientMock(controller)
				require.NotNil(t, aggregatorClientMock)

				aggregatorClientMock.
					GetTimeSummaryMock.
					Times(1).
					Expect(minimock.AnyContext, "", "2025-05-06").
					Return(nil, errors.New("internal error"))

				return aggregatorClientMock
			},
			APIClient: func() APIClient {
				apiClientMock := NewAPIClientMock(controller)
				require.NotNil(t, apiClientMock)

				return apiClientMock
			},
			date: "2025-05-06",

			expectedErr: fmt.Errorf("error getting time summary: %w", errors.New("internal error")),
		},
		"empty time summary by all users": {
			notifier: func() notifier {
				notifierMock := NewNotifierMock(controller)
				require.NotNil(t, notifierMock)

				return notifierMock
			},
			aggregatorClient: func() aggregatorClient {
				aggregatorClientMock := NewAggregatorClientMock(controller)
				require.NotNil(t, aggregatorClientMock)

				aggregatorClientMock.
					GetTimeSummaryMock.
					Times(1).
					Expect(minimock.AnyContext, "", "2025-05-06").
					Return(&pb.GetSummaryResponse{Summary: []*pb.Summary{}}, nil)

				return aggregatorClientMock
			},
			APIClient: func() APIClient {
				apiClientMock := NewAPIClientMock(controller)
				require.NotNil(t, apiClientMock)

				return apiClientMock
			},
			date: "2025-05-06",

			expectedErr: nil,
		},
		"success send notify": {
			notifier: func() notifier {
				notifierMock := NewNotifierMock(controller)
				require.NotNil(t, notifierMock)

				notifierMock.
					SendPreviousDayInfoMessageMock.
					Expect(minimock.AnyContext, telegram.ParamsPreviousDayInfo{
						TelegramID:   12312312,
						SecondsStart: helper.SecondsToTime(1746548943),
						SecondsEnd:   helper.SecondsToTime(1746548950),
						Hours:        34,
						Minutes:      12,
						Breaks:       "",
					}).
					Return(nil)

				return notifierMock
			},
			aggregatorClient: func() aggregatorClient {
				aggregatorClientMock := NewAggregatorClientMock(controller)
				require.NotNil(t, aggregatorClientMock)

				aggregatorClientMock.
					GetTimeSummaryMock.
					Times(1).
					Expect(minimock.AnyContext, "", "2025-05-06").
					Return(&pb.GetSummaryResponse{Summary: expectedTimeSummaryByUsers}, nil)

				return aggregatorClientMock
			},
			APIClient: func() APIClient {
				apiClientMock := NewAPIClientMock(controller)
				require.NotNil(t, apiClientMock)

				apiClientMock.
					GetUserByMacAddressMock.
					Times(1).
					Expect(minimock.AnyContext, "00:1b:63:84:45:e6").
					Return(&pbapiv1.GetUserByMacAddressResponse{
						User: &pbapiv1.User{
							Id:         1,
							Name:       "Ivan",
							Surname:    "Ivanov",
							Lastname:   "Ivanovich",
							BirthDate:  "2020-03-01",
							Email:      "ivan@gmail.com",
							MacAddress: "00:1b:63:84:45:e6",
							TelegramId: 12312312,
							Role:       "employee",
							Department: 1,
							Position:   "programmer",
						},
					}, nil)

				return apiClientMock
			},
			date: "2025-05-06",

			expectedErr: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			previousDayInfoService, err := NewPreviousDayInfo(
				test.APIClient(),
				test.aggregatorClient(),
				test.notifier(),
				clockwork.NewFakeClock(),
				12,
				0,
				0,
			)
			require.Nil(t, err)
			require.NotNil(t, previousDayInfoService)

			err = previousDayInfoService.sendAllUsersNotify(minimock.AnyContext, test.date)
			if test.expectedErr != nil {
				assert.Error(t, test.expectedErr, err)
			}
		})
	}
}
