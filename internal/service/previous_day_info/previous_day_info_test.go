package previous_day_info

import (
	"context"
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

		/*notifierMock.
		SendPreviousDayInfoMessageMock.
		Expect(minimock.AnyContext, telegram.ParamsPreviousDayInfo{}).
		Return(nil)*/

		//const macAddress string = "00:1b:63:84:45:e6"

		aggregatorClientMock := NewAggregatorClientMock(controller)
		require.NotNil(t, aggregatorClientMock)
		/*aggregatorClientMock.
		GetTimeSummaryMock.
		Expect(minimock.AnyContext, macAddress, "2025-05-06").
		Return(&pb.GetSummaryResponse{}, nil)*/

		apiClientMock := NewAPIClientMock(controller)
		require.NotNil(t, apiClientMock)
		/*apiClientMock.
		GetUserByMacAddressMock.
		Expect(minimock.AnyContext, macAddress).
		Return(&pbapiv1.GetUserByMacAddressResponse{}, nil)*/

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

		handler := func(_ context.Context) error {
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

	t.Run("context cancelled", func(t *testing.T) {
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

		handler := func(_ context.Context) error {
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
