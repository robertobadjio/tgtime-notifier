package command

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/robertobadjio/tgtime-aggregator/pkg/api/time_v1"
	"github.com/robertobadjio/tgtime-notifier/internal/aggregator"
	"github.com/robertobadjio/tgtime-notifier/internal/api_pb"
	"github.com/robertobadjio/tgtime-notifier/internal/config"
)

// WorkingTimeCommand Команда "Рабочее время"
type WorkingTimeCommand struct {
	TelegramID int64
}

// GetMessage Метод получения текста команды
func (wtc WorkingTimeCommand) GetMessage(ctx context.Context) (string, error) {
	tgTimeAggregatorConfig, err := config.NewTgTimeAggregatorConfig()
	if err != nil {
		return "", fmt.Errorf("error loading config: %w", err)
	}

	tgTimeAPIConfig, err := config.NewTgTimeAPIConfig()
	if err != nil {
		return "", fmt.Errorf("error loading config: %w", err)
	}

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	aggregatorClient := aggregator.NewClient(tgTimeAggregatorConfig, logger)
	apiClient := api_pb.NewClient(tgTimeAPIConfig, logger)
	user, err := apiClient.GetUserByTelegramID(ctx, wtc.TelegramID)
	if err != nil {
		return "", fmt.Errorf("error getting user by telegram id: %w", err)
	}

	timeSummaryResponse, err := aggregatorClient.GetTimeSummary(
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
