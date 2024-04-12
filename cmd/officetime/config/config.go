package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"regexp"
)

type Config struct {
	BotToken          string
	RouterAddress     string
	RouterUserName    string
	RouterPassword    string
	WebHookPath       string
	ApiURL            string
	ApiMasterEmail    string
	ApiMasterPassword string
}

const projectDirName = "cloud-time-tracker"

func init() {
	loadEnv()
}

func New() *Config {
	return &Config{
		BotToken:          getEnv("BOT_TOKEN", ""),
		RouterAddress:     getEnv("ROUTER_ADDRESS", ""),
		RouterUserName:    getEnv("ROUTER_USER_NAME", ""),
		RouterPassword:    getEnv("ROUTER_PASSWORD", ""),
		WebHookPath:       getEnv("WEBHOOK_PATH", ""),
		ApiURL:            getEnv("API_URL", ""),
		ApiMasterEmail:    getEnv("API_MASTER_EMAIL", ""),
		ApiMasterPassword: getEnv("API_MASTER_PASSWORD", ""),
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
