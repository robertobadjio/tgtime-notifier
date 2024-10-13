package command

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/log"
	"os"
	"strings"
	"tgtime-notifier/internal/aggregator"
	"tgtime-notifier/internal/api_pb"
	"tgtime-notifier/internal/config"
	"time"
)

type WorkingTimeCommand struct {
	TelegramId int64
}

type Break struct {
	BeginTime int64 `json:"beginTime"` // TODO: rename StartTime
	EndTime   int64 `json:"endTime"`
}

func (wtc WorkingTimeCommand) GetMessage(ctx context.Context) (string, error) {
	cfg := config.New()

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	aggregatorClient := aggregator.NewClient(*cfg, logger)
	apiClient := api_pb.NewClient(*cfg, logger)
	user, err := apiClient.GetUserByTelegramId(ctx, wtc.TelegramId)
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
	if len(timeSummaryResponse.TimeSummary) == 0 {
		return "", fmt.Errorf("time summary not found mac address " + user.User.MacAddress + " date " + getNow().Format("2006-01-02"))
	}

	if timeSummaryResponse.TimeSummary[0].SecondsStart == 0 {
		return "Вы сегодня не были в офисе", nil
	} else {
		hours, minutes := secondsToHM(int(timeSummaryResponse.TimeSummary[0].Seconds))
		beginTime := time.Unix(timeSummaryResponse.TimeSummary[0].SecondsStart, 0)
		mes := fmt.Sprintf(
			"Сегодня Вы в офисе с %s\nУчтенное время %d ч. %d м.",
			beginTime.Format("15:04"),
			hours,
			minutes,
		)

		var breaksRaw []*Break

		// TODO: По GRPC отдавать сразу срез
		_ = json.Unmarshal([]byte(timeSummaryResponse.TimeSummary[0].GetBreaksJson()), &breaksRaw)
		breaks := breaksToString(buildBreaks(breaksRaw))
		if breaks != "" {
			mes += fmt.Sprintf("\nПерерывы %s", breaks)
		}

		return mes, nil
	}
}

func secondsToHM(seconds int) (int, int) {
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

func buildBreaks(breaks []*Break) []string {
	var output []string
	for _, item := range breaks {
		beginTime := time.Unix(item.BeginTime, 0)
		endTime := time.Unix(item.EndTime, 0)
		output = append(
			output,
			fmt.Sprintf("%s - %s", beginTime.Format("15:04"), endTime.Format("15:04")))
	}

	return output
}

func breaksToString(breaks []string) string {
	return strings.Join(breaks, ", ")
}
