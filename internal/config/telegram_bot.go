package config

import (
	"fmt"
)

const (
	botTokenEnvParam    = "BOT_TOKEN"
	webhookPathEnvParam = "WEBHOOK_PATH"
	webhookLinkEnvParam = "WEBHOOK_LINK"
)

// TelegramBotConfig ...
type TelegramBotConfig struct {
	token       string
	webhookPath string
	webhookLink string
}

// NewTelegramBotConfig ...
func NewTelegramBotConfig(os OS) (*TelegramBotConfig, error) {
	if os == nil {
		return nil, fmt.Errorf("os must not be nil")
	}

	token := os.Getenv(botTokenEnvParam)
	if len(token) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", botTokenEnvParam)
	}

	webhookPath := os.Getenv(webhookPathEnvParam)
	webhookLink := os.Getenv(webhookLinkEnvParam)

	return &TelegramBotConfig{
		token:       token,
		webhookPath: webhookPath,
		webhookLink: webhookLink,
	}, nil
}

// Token ...
func (t TelegramBotConfig) Token() string {
	return t.token
}

// WebhookPath ...
func (t TelegramBotConfig) WebhookPath() string {
	return t.webhookPath
}

// WebhookLink ...
func (t TelegramBotConfig) WebhookLink() string {
	return t.webhookLink
}
