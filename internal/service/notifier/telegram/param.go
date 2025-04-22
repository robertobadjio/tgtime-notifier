package telegram

import (
	"time"

	TGBotAPI "github.com/go-telegram-bot-api/telegram-bot-api"
)

// ParamsWelcomeMessage ...
type ParamsWelcomeMessage struct {
	TelegramID int64
}

// ParamsUpdate ...
type ParamsUpdate struct {
	Update TGBotAPI.Update
}

// ParamsWorkingTime ...
type ParamsWorkingTime struct {
}

// ParamsPreviousDayInfo ...
type ParamsPreviousDayInfo struct {
	TelegramID   int64
	SecondsStart time.Time
	SecondsEnd   time.Time
	Hours        int64
	Minutes      int64
	Breaks       string
}
