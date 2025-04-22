package helper

import (
	"fmt"
	"strings"
	"time"

	"github.com/robertobadjio/tgtime-aggregator/pkg/api/time_v1"
)

// SecondsToHM ...
func SecondsToHM(seconds int64) (int64, int64) {
	hours := seconds / 3600
	minutes := (seconds / 60) - (hours * 60)

	return hours, minutes
}

// GetNow ...
func GetNow() time.Time {
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

// BreaksToString ...
func BreaksToString(breaks []string) string {
	return strings.Join(breaks, ", ")
}

// SecondsToTime ...
func SecondsToTime(seconds int64) time.Time {
	return time.Unix(seconds, 0)
}
