package config

import (
	"errors"
	"fmt"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestNewPromConfig(t *testing.T) {
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
				osMock.GetenvMock.Expect(promAppPortEnvName).Times(1).Return("2112")

				return osMock
			},
			expectedNilObj: false,
			expectedErr:    nil,
		},
		"create config with empty port": {
			os: func() OS {
				osMock := NewOSMock(mc)
				osMock.GetenvMock.Expect(promAppPortEnvName).Times(1).Return("")

				return osMock
			},
			expectedNilObj: true,
			expectedErr:    fmt.Errorf("environment variable %s not set", promAppPortEnvName),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cfg, err := NewPromConfig(test.os())
			assert.Equal(t, test.expectedErr, err)

			if test.expectedNilObj {
				assert.Nil(t, cfg)
			} else {
				assert.NotNil(t, cfg)
			}
		})
	}
}

func TestPromConfig_Address(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)

	tests := map[string]struct {
		os func() OS

		expectedAddress string
	}{
		"get address": {
			os: func() OS {
				osMock := NewOSMock(mc)
				osMock.GetenvMock.Expect(promAppPortEnvName).Times(1).Return("2112")

				return osMock
			},
			expectedAddress: ":2112",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cfg, err := NewPromConfig(test.os())
			assert.Nil(t, err)
			assert.NotNil(t, cfg)
			assert.Equal(t, test.expectedAddress, cfg.Address())
		})
	}
}

func TestPromConfig_Path(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)

	tests := map[string]struct {
		os func() OS

		expectedPath string
	}{
		"get address": {
			os: func() OS {
				osMock := NewOSMock(mc)
				osMock.GetenvMock.Expect(promAppPortEnvName).Times(1).Return("2112")

				return osMock
			},
			expectedPath: "/metrics",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cfg, err := NewPromConfig(test.os())
			assert.Nil(t, err)
			assert.NotNil(t, cfg)
			assert.Equal(t, test.expectedPath, cfg.Path())
		})
	}
}
