package config

import (
	"log"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

// Config Конфига приложения
type Config struct {
	BotToken             string
	RouterAddress        string
	RouterUserName       string
	RouterPassword       string
	WebHookPath          string
	WebHookLink          string
	APIURL               string
	APIMasterEmail       string
	APIMasterPassword    string
	KafkaHost            string
	KafkaPort            string
	TgTimeAggregatorHost string
	TgTimeAggregatorPort string
	TgTimeAPIHost        string
	TgTimeAPIPort        string
}

const projectDirName = "tgtime-notifier"

func init() {
	loadEnv()
}

// New Конструктор конфига приложения
func New() *Config {
	return &Config{
		BotToken:             getEnv("BOT_TOKEN", ""),
		RouterAddress:        getEnv("ROUTER_ADDRESS", ""),
		RouterUserName:       getEnv("ROUTER_USER_NAME", ""),
		RouterPassword:       getEnv("ROUTER_PASSWORD", ""),
		WebHookPath:          getEnv("WEBHOOK_PATH", ""),
		WebHookLink:          getEnv("WEBHOOK_LINK", ""),
		APIURL:               getEnv("API_URL", ""),
		APIMasterEmail:       getEnv("API_MASTER_EMAIL", ""),
		APIMasterPassword:    getEnv("API_MASTER_PASSWORD", ""),
		KafkaHost:            getEnv("KAFKA_HOST", ""),
		KafkaPort:            getEnv("KAFKA_PORT", ""),
		TgTimeAggregatorHost: getEnv("TGTIME_AGGREGATOR_HOST", ""),
		TgTimeAggregatorPort: getEnv("TGTIME_AGGREGATOR_PORT", ""),
		TgTimeAPIHost:        getEnv("TGTIME_API_HOST", ""),
		TgTimeAPIPort:        getEnv("TGTIME_API_PORT", ""),
	}
}

func loadEnv() {
	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	err := godotenv.Load(string(rootPath) + `/.env`)
	if err != nil {
		log.Fatal("Problem loading .env file")
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
