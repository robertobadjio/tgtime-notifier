package config

import (
	"errors"
	"fmt"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestTgTimeAPIConfig_New(t *testing.T) {
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
		"create config with port": {
			os: func() OS {
				osMock := NewOSMock(mc)
				osMock.GetenvMock.When(tgTimeAPIHost).Then("127.0.0.1")
				osMock.GetenvMock.When(tgTimeAPIPort).Then("8080")

				return osMock
			},
			expectedNilObj: false,
			expectedErr:    nil,
		},
		"create config with empty port": {
			os: func() OS {
				osMock := NewOSMock(mc)
				osMock.GetenvMock.When(tgTimeAPIHost).Then("127.0.0.1")
				osMock.GetenvMock.When(tgTimeAPIPort).Then("")

				return osMock
			},
			expectedNilObj: true,
			expectedErr:    fmt.Errorf("environment variable %s must be set", tgTimeAPIPort),
		},
		"create config with empty host": {
			os: func() OS {
				osMock := NewOSMock(mc)
				osMock.GetenvMock.When(tgTimeAPIHost).Then("")
				osMock.GetenvMock.When(tgTimeAPIPort).Then("8080")

				return osMock
			},
			expectedNilObj: false,
			expectedErr:    nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cfg, err := NewTgTimeAPIConfig(test.os())
			assert.Equal(t, test.expectedErr, err)

			if test.expectedNilObj {
				assert.Nil(t, cfg)
			} else {
				assert.NotNil(t, cfg)
			}
		})
	}
}

func TestTgTimeAPIConfig_Address(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)

	tests := map[string]struct {
		os func() OS

		expectedAddress string
	}{
		"get address with host": {
			os: func() OS {
				osMock := NewOSMock(mc)
				osMock.GetenvMock.When(tgTimeAPIHost).Then("127.0.0.1")
				osMock.GetenvMock.When(tgTimeAPIPort).Then("8080")

				return osMock
			},
			expectedAddress: "127.0.0.1:8080",
		},
		"get address without host": {
			os: func() OS {
				osMock := NewOSMock(mc)
				osMock.GetenvMock.When(tgTimeAPIHost).Then("")
				osMock.GetenvMock.When(tgTimeAPIPort).Then("8080")

				return osMock
			},
			expectedAddress: ":8080",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cfg, err := NewTgTimeAPIConfig(test.os())
			assert.Nil(t, err)
			assert.NotNil(t, cfg)
			assert.Equal(t, test.expectedAddress, cfg.Address())
		})
	}
}
