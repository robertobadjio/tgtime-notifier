package config

import (
	"fmt"
	"strconv"
)

const (
	hourEnvParam = "PREVIOUS_DAY_INFO_HOUR"
)

const defaultHour int = 12

type PreviousDayInfoConfig struct {
	hour int
}

// NewPreviousDayInfoConfig ...
func NewPreviousDayInfoConfig(os OS) (*PreviousDayInfoConfig, error) {
	if os == nil {
		return nil, fmt.Errorf("os must not be nil")
	}

	hourRaw := os.Getenv(hourEnvParam)
	var hour int
	if len(hourRaw) == 0 {
		hour = defaultHour
	} else {
		hourInt, err := strconv.Atoi(hourRaw)
		if err != nil {
			return nil, fmt.Errorf("error parsing %s env var: %s", hourEnvParam, err)
		}
		hour = hourInt
	}

	return &PreviousDayInfoConfig{
		hour: hour,
	}, nil
}

func (pdi PreviousDayInfoConfig) Hour() int {
	return pdi.hour
}
