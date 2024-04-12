package telegram

import (
	"cloud-time-tracker/cmd/officetime/api"
	"cloud-time-tracker/cmd/officetime/service"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
	"time"
)

const (
	buttonWorkingTime              = "⏳ Рабочее время"
	buttonStatCurrentWorkingPeriod = "🗓 Статистика за рабочий период"
	buttonStart                    = "/start"
)

type User struct {
	TelegramId int64
	UserId     int64
}

var Users map[string]User

type Telegram struct {
	users map[string]User
	bot   *tgbotapi.BotAPI
}

func NewTelegram(bot *tgbotapi.BotAPI, users map[string]User) *Telegram {
	return &Telegram{bot: bot, users: users}
}

func (t Telegram) Info(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	macAddress := t.searchMacAddressByTelegramId(int64(update.Message.From.ID))

	if macAddress == "" {
		panic("Mac address not found")
	}

	// TODO: стратегия
	if update.Message.Text == buttonWorkingTime {
		c := api.NewOfficeTimeClient()
		//response, err := c.GetTimesByDate(int(UserUserId[macAddress]), time.Now())
		response, err := c.GetTimesByDate(int(t.users[macAddress].UserId), time.Now())
		if err != nil {
			panic(err)
		}

		var messageTelegram tgbotapi.MessageConfig
		if response.BeginTime == 0 {
			messageTelegram = tgbotapi.NewMessage(int64(update.Message.From.ID), "Вы сегодня не были в офисе")
		} else {
			hours, minutes := secondsToHM(response.Total)
			beginTime := time.Unix(response.BeginTime, 0)
			message := fmt.Sprintf("Сегодня Вы в офисе с %s\nУчтенное время %d ч. %d м.", beginTime.Format("15:04"), hours, minutes)
			breaks := breaksToString(buildBreaks(response.Breaks))
			if breaks != "" {
				message += fmt.Sprintf("\nПерерывы %s", breaks)
			}
			messageTelegram = tgbotapi.NewMessage(int64(update.Message.From.ID), message)
		}
		bot.Send(setKeyboard(messageTelegram))
	} else if update.Message.Text == buttonStatCurrentWorkingPeriod {
		c := api.NewOfficeTimeClient()
		//result, err := c.GetStatByWorkingPeriod(int(UserUserId[macAddress]), service.GetCurrentPeriod())
		result, err := c.GetStatByWorkingPeriod(int(t.users[macAddress].UserId), service.GetCurrentPeriod())
		if err != nil {
			panic(err)
		}

		message := tgbotapi.NewMessage(
			int64(update.Message.From.ID),
			fmt.Sprintf("Статистика за период с %s до %s\nВсего в этом месяце %d из %d часов", result.StartWorkingDate, result.EndWorkingDate, result.WorkingHours, result.TotalWorkingHours))

		bot.Send(setKeyboard(message))
	} else if update.Message.Text == buttonStart {
		message := tgbotapi.NewMessage(
			int64(update.Message.From.ID),
			"Добро пожаловать. Используйте кнопки для получения информации")
		bot.Send(setKeyboard(message))
	} else {
		message := tgbotapi.NewMessage(
			int64(update.Message.From.ID),
			"Неизвестная команда")
		bot.Send(setKeyboard(message))
	}
}

func (t Telegram) SendMessage(bot *tgbotapi.BotAPI, telegramId int64) {
	msg := tgbotapi.NewMessage(telegramId, "Вы пришли в офис")
	bot.Send(msg)
}

func (t Telegram) searchMacAddressByTelegramId(telegramId int64) string {
	for macAddress, value := range t.users {
		if value.TelegramId == telegramId {
			return macAddress
		}
	}

	return ""
}

// Keyboard
func setKeyboard(message tgbotapi.MessageConfig) tgbotapi.MessageConfig {
	message.ReplyMarkup = tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(buttonWorkingTime), tgbotapi.NewKeyboardButton(buttonStatCurrentWorkingPeriod)))
	return message
}

// TODO: Перенести в микросервис aggregator
func previousDayInfo(bot *tgbotapi.BotAPI) {
	moscow, _ := time.LoadLocation("Europe/Moscow")
	yesterday := time.Now().In(moscow).Add(-24 * time.Hour)

	c := api.NewOfficeTimeClient()
	result, err := c.GetAllUsers()
	if err != nil {
		panic(err)
	}

	users := result.Users

	for _, user := range users {
		var hours, minutes int
		var beginTimeSeconds, endTimeSeconds int64
		c := api.NewOfficeTimeClient()
		response, err := c.GetTimesByPeriod(user.Id, service.GetCurrentPeriod())
		if err != nil {
			panic(err)
		}
		for _, timeResponse := range response.Times {
			if timeResponse.Date != yesterday.Format("2006-01-02") {
				continue
			}

			if timeResponse.Total == 0 {
				break
			}

			hours, minutes = secondsToHM(timeResponse.Total)
			beginTimeSeconds = timeResponse.BeginTime
			endTimeSeconds = timeResponse.EndTime

			beginTime := time.Unix(beginTimeSeconds, 0)
			endTime := time.Unix(endTimeSeconds, 0)

			breaks := breaksToString(buildBreaks(timeResponse.Breaks))
			message := fmt.Sprintf("Вчера Вы были в офисе с %s до %s\nУчтенное время %d ч. %d м.", beginTime.Format("15:04"), endTime.Format("15:04"), hours, minutes)
			if "" != breaks {
				message += fmt.Sprintf("\nПерерывы %s\n", breaks)
			}
			bot.Send(setKeyboard(tgbotapi.NewMessage(user.TelegramId, message)))
		}
	}
}

func breaksToString(breaks []string) string {
	return strings.Join(breaks, ", ")
}

func buildBreaks(breaks []api.BreakResponse) []string {
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

func secondsToHM(seconds int) (int, int) {
	hours := seconds / 3600
	minutes := (seconds / 60) - (hours * 60)

	return hours, minutes
}
