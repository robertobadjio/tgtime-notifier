package main

import (
	"cloud-time-tracker/cmd/officetime/api"
	"cloud-time-tracker/cmd/officetime/config"
	"cloud-time-tracker/cmd/officetime/service"
	"cloud-time-tracker/cmd/officetime/telegram"
	"cloud-time-tracker/internal/background"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

//var homes map[string]map[string]bool
//var telegramBot *tgbotapi.BotAPI

/*func init() {
	homes = make(map[string]map[string]bool)
}*/

func main() {
	telegramBot, _ := initTelegramBot()
	users, _ := getUsers()

	t := telegram.NewTelegram(telegramBot, users)

	// TODO: Сервис упал, нужно откуда-то восстановить данные
	/*userStore := user_store.NewStore(*users)

	// TODO:
	// Нужно ходить в time tracker и получать только что подключившихся пользователей
	// Можно сделать через pub/sub, redis или ходить синхронно каждые 10 секунд
	f := func() {
		// TODO: Получить подключившегося пользователя
		//t.SendMessage(telegramBot, telegram.Users[macAddress].TelegramId)
	}
	bc := background.NewBackground(time.Duration(10)*time.Second, f)
	bc.Start()*/

	// TODO: В 12 дня посылать информацио о предыдущем дне
	f2 := func() {
		year, month, day := service.GetNow().Date()
		// TODO: в параметры
		dateTimeBegin := time.Date(year, month, day, 12, 00, 00, 0, service.GetNow().Location())
		dateTimeEnd := time.Date(year, month, day, 12, 00, 30, 0, service.GetNow().Location())

		if service.GetNow().Before(dateTimeBegin) || service.GetNow().After(dateTimeEnd) {
			return
		}
		// TODO: previousDayInfo
		// Получить данные из микросервиса aggregator
		// Реализовать эндпоинт получения timeSummary по всем пользователям
	}
	bc2 := background.NewBackground(time.Duration(60)*time.Second, f2)
	bc2.Start()
	//err := initHomes()
	/*if err != nil {
		go eachEvery30Seconds()
	}*/

	cfg := config.New()
	updates := telegramBot.ListenForWebhook("/" + cfg.WebHookPath)
	go http.ListenAndServe(":8441", nil)
	for update := range updates {
		t.Info(update, telegramBot)
	}
}

/*func eachEvery30Seconds() {
	for {
		c := api.NewOfficeTimeClient()
		result, err := c.GetAllRouters()
		if err != nil {
			fmt.Println("Ошибка при получении списка роутеров")
		} else {
			for _, routerDevice := range result.Routers {
				var macAddresses []string
				macAddresses = router.Router(routerDevice)

				for _, macAddress := range macAddresses { // TODO: Отправлять пачкой для каждого роутера
					//addFalseToHomesIfNotFound(macAddress)
					//sendHomeMessage(strings.ToLower(macAddress))
				}
			}

			telegram.PreviousDayInfo(telegramBot)
		}

		time.Sleep(30000 * time.Millisecond)
	}
}*/

func initTelegramBot() (*tgbotapi.BotAPI, error) {
	cfg := config.New()
	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, err
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook("https://demo.tgtime.ru/telegram")) // TODO: в параметры
	if err != nil {
		return nil, err
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		return nil, err
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	return bot, nil
}

func getUsers() (*api.Users, error) {
	c := api.NewOfficeTimeClient()
	result, err := c.GetAllUsers()
	if err != nil {
		return nil, err
	}

	return result, nil
}

/*func initHomes() error {
	c := api.NewOfficeTimeClient()
	result, err := c.GetAllUsers()
	if err != nil {
		panic(err)
	}

	users := result.Users
	homes[getDateNow()] = make(map[string]bool)
	for _, user := range users {
		c := api.NewOfficeTimeClient()
		times, err := c.GetTimesByPeriod(user.Id, service.GetCurrentPeriod())
		if err != nil {
			return fmt.Errorf("Нет периодов")
		}
		for _, timeStruct := range times.Times {
			if timeStruct.Date != getDateNow() {
				continue
			}
			homes[getDateNow()][user.MacAddress] = timeStruct.BeginTime > 0
		}
	}
	return nil
}*/

/*func addFalseToHomesIfNotFound(macAddress string) {
	_, found := homes[getDateNow()]
	if !found {
		homes[getDateNow()] = make(map[string]bool)
	}

	homes[getDateNow()][macAddress] = false
}*/

/*func sendHomeMessage(macAddress string) {
	if homes[getDateNow()][macAddress] == false {
		_, found := telegram.Users[macAddress]
		if found {
			telegram.SendMessage(telegramBot, telegram.Users[macAddress].TelegramId)
		}
		homes[getDateNow()][macAddress] = true
	}
}*/

/*func getDateNow() string {
	return service.GetNow().Format("2006-01-02")
}*/
