package config

import (
	"errors"
	"fmt"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestTelegramBotConfig_New(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)

	tests := map[string]struct {
		os func() OS

		expectedNilObj bool
		expectedErr    error
	}{
		"create config with empty OS": {
			os: func() OS {
				return nil
			},
			expectedNilObj: true,
			expectedErr:    errors.New("os must not be nil"),
		},
		"create config with empty token": {
			os: func() OS {
				osMock := NewOSMock(mc)
				osMock.GetenvMock.Expect(botTokenEnvParam).Times(1).Return("")

				return osMock
			},
			expectedNilObj: true,
			expectedErr:    fmt.Errorf("environment variable %s must be set", botTokenEnvParam),
		},
		"create config with empty webhook path": {
			os: func() OS {
				osMock := NewOSMock(mc)
				osMock.GetenvMock.When(botTokenEnvParam).Then("FHD3DO")
				osMock.GetenvMock.When(webhookPathEnvParam).Then("")

				return osMock
			},
			expectedNilObj: true,
			expectedErr:    fmt.Errorf("environment variable %s must be set", webhookPathEnvParam),
		},
		"create config with empty webhook link": {
			os: func() OS {
				osMock := NewOSMock(mc)
				osMock.GetenvMock.When(botTokenEnvParam).Then("2343434:sad7Dsad_DSADDk3k4")
				osMock.GetenvMock.When(webhookPathEnvParam).Then("telegram")
				osMock.GetenvMock.When(webhookLinkEnvParam).Then("")

				return osMock
			},
			expectedNilObj: true,
			expectedErr:    fmt.Errorf("environment variable %s must be set", webhookLinkEnvParam),
		},
		"create config": {
			os: func() OS {
				osMock := NewOSMock(mc)
				osMock.GetenvMock.When(botTokenEnvParam).Then("2343434:sad7Dsad_DSADDk3k4")
				osMock.GetenvMock.When(webhookPathEnvParam).Then("telegram")
				osMock.GetenvMock.When(webhookLinkEnvParam).Then("https://example.ru/telegram")

				return osMock
			},
			expectedNilObj: false,
			expectedErr:    nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cfg, err := NewTelegramBotConfig(test.os())
			assert.Equal(t, test.expectedErr, err)

			if test.expectedNilObj {
				assert.Nil(t, cfg)
			} else {
				assert.NotNil(t, cfg)
			}
		})
	}
}

func TestTelegramBotConfig_GetToken(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)

	tests := map[string]struct {
		os func() OS

		expectedToken string
	}{
		"get token": {
			os: func() OS {
				osMock := NewOSMock(mc)
				osMock.GetenvMock.When(botTokenEnvParam).Then("2343434:sad7Dsad_DSADDk3k4")
				osMock.GetenvMock.When(webhookPathEnvParam).Then("telegram")
				osMock.GetenvMock.When(webhookLinkEnvParam).Then("https://example.ru/telegram")

				return osMock
			},
			expectedToken: "2343434:sad7Dsad_DSADDk3k4",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cfg, err := NewTelegramBotConfig(test.os())
			assert.Nil(t, err)
			assert.NotNil(t, cfg)
			assert.Equal(t, test.expectedToken, cfg.Token())
		})
	}
}

func TestTelegramBotConfig_WebhookPath(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)

	tests := map[string]struct {
		os func() OS

		expectedPath string
	}{
		"get token": {
			os: func() OS {
				osMock := NewOSMock(mc)
				osMock.GetenvMock.When(botTokenEnvParam).Then("2343434:sad7Dsad_DSADDk3k4")
				osMock.GetenvMock.When(webhookPathEnvParam).Then("telegram")
				osMock.GetenvMock.When(webhookLinkEnvParam).Then("https://example.ru/telegram")

				return osMock
			},
			expectedPath: "telegram",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cfg, err := NewTelegramBotConfig(test.os())
			assert.Nil(t, err)
			assert.NotNil(t, cfg)
			assert.Equal(t, test.expectedPath, cfg.WebhookPath())
		})
	}
}

func TestTelegramBotConfig_WebhookLink(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)

	tests := map[string]struct {
		os func() OS

		expectedLink string
	}{
		"get token": {
			os: func() OS {
				osMock := NewOSMock(mc)
				osMock.GetenvMock.When(botTokenEnvParam).Then("2343434:sad7Dsad_DSADDk3k4")
				osMock.GetenvMock.When(webhookPathEnvParam).Then("telegram")
				osMock.GetenvMock.When(webhookLinkEnvParam).Then("https://example.ru/telegram")

				return osMock
			},
			expectedLink: "https://example.ru/telegram",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cfg, err := NewTelegramBotConfig(test.os())
			assert.Nil(t, err)
			assert.NotNil(t, cfg)
			assert.Equal(t, test.expectedLink, cfg.WebhookLink())
		})
	}
}
