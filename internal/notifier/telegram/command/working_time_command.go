package command

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/robertobadjio/tgtime-aggregator/pkg/api/time_v1"

	"github.com/robertobadjio/tgtime-notifier/internal/aggregator"
	"github.com/robertobadjio/tgtime-notifier/internal/api_pb"
)

type workingTimeCommand struct {
	TGTimeAPIClient        api_pb.Client
	TGTimeAggregatorClient aggregator.Client
	TelegramID             int64
}

// NewWorkingTimeCommand ???
func NewWorkingTimeCommand(TGTimeAPIClient api_pb.Client, TGTimeAggregatorClient aggregator.Client, telegramID int64) Command {
	return &workingTimeCommand{
		TGTimeAPIClient:        TGTimeAPIClient,
		TGTimeAggregatorClient: TGTimeAggregatorClient,
		TelegramID:             telegramID,
	}
}

// GetMessage ???
func (wtc *workingTimeCommand) GetMessage(ctx context.Context) (string, error) {
	user, err := wtc.TGTimeAPIClient.GetUserByTelegramID(ctx, wtc.TelegramID)
	if err != nil {
		return "", fmt.Errorf("error getting user by telegram id: %w", err)
	}

	timeSummaryResponse, err := wtc.TGTimeAggregatorClient.GetTimeSummary(
		ctx,
		user.User.MacAddress,
		getNow().Format("2006-01-02"),
	)
	if err != nil {
		return "", fmt.Errorf("error getting time summary: %w", err)
	}
	if len(timeSummaryResponse.Summary) == 0 {
		return "", fmt.Errorf("time summary not found")
	}

	if timeSummaryResponse.Summary[0].SecondsStart == 0 {
		return "Вы сегодня не были в офисе", nil
	}

	hours, minutes := SecondsToHM(timeSummaryResponse.Summary[0].Seconds)
	beginTime := SecondsToTime(timeSummaryResponse.Summary[0].SecondsStart)
	mes := fmt.Sprintf(
		"Сегодня Вы в офисе с %s\nУчтенное время %d ч. %d м.",
		beginTime.Format("15:04"),
		hours,
		minutes,
	)

	breaks := BreaksToString(BuildBreaks(timeSummaryResponse.Summary[0].GetBreaks()))
	if breaks != "" {
		mes += fmt.Sprintf("\nПерерывы %s", breaks)
	}

	return mes, nil
}

// SecondsToHM ???
func SecondsToHM(seconds int64) (int64, int64) {
	hours := seconds / 3600
	minutes := (seconds / 60) - (hours * 60)

	return hours, minutes
}

func getNow() time.Time {
	return time.Now().In(getMoscowLocation())
}

func getMoscowLocation() *time.Location {
	moscowLocation, _ := time.LoadLocation("Europe/Moscow")
	return moscowLocation
}

// BuildBreaks ???
func BuildBreaks(breaks []*time_v1.Break) []string {
	var output []string
	for _, item := range breaks {
		beginTime := time.Unix(item.SecondsStart, 0)
		endTime := time.Unix(item.SecondsEnd, 0)
		output = append(
			output,
			fmt.Sprintf("%s - %s", beginTime.Format("15:04"), endTime.Format("15:04")))
	}

	return output
}

// BreaksToString ???
func BreaksToString(breaks []string) string {
	return strings.Join(breaks, ", ")
}

// SecondsToTime ???
func SecondsToTime(seconds int64) time.Time {
	return time.Unix(seconds, 0)
}
