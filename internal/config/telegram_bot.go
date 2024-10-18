package config

import (
	"fmt"
	"os"
)

const botTokenEnvParam = "BOT_TOKEN"
const webhookPathEnvParam = "WEBHOOK_PATH"
const webhookLinkEnvParam = "WEBHOOK_LINK"

// TelegramBotConfig ???
type TelegramBotConfig interface {
	GetToken() string
	GetWebhookPath() string
	GetWebhookLink() string
}

type telegramBotConfig struct {
	token       string
	webhookPath string
	webhookLink string
}

// GetToken ???
func (t *telegramBotConfig) GetToken() string {
	return t.token
}

// GetWebhookPath ???
func (t *telegramBotConfig) GetWebhookPath() string {
	return t.webhookPath
}

// GetWebhookLink ???
func (t *telegramBotConfig) GetWebhookLink() string {
	return t.webhookLink
}

// NewTelegramBotConfig ???
func NewTelegramBotConfig() (TelegramBotConfig, error) {
	token := os.Getenv(botTokenEnvParam)
	if len(token) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", botTokenEnvParam)
	}

	webhookPath := os.Getenv(webhookPathEnvParam)
	if len(token) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", webhookPathEnvParam)
	}

	webhookLink := os.Getenv(webhookLinkEnvParam)
	if len(token) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", webhookLinkEnvParam)
	}

	return &telegramBotConfig{
		token:       token,
		webhookPath: webhookPath,
		webhookLink: webhookLink,
	}, nil
}
