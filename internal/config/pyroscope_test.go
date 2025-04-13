package config

import (
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestPyroscopeConfig_New(t *testing.T) {
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
		"create config with host and port": {
			os: func() OS {
				osMock := NewOSMock(mc)
				osMock.GetenvMock.When(pyroscopeHostEnvName).Then("127.0.0.1")
				osMock.GetenvMock.When(pyroscopePortEnvName).Then("4040")

				return osMock
			},
			expectedNilObj: false,
			expectedErr:    nil,
		},
		"create config with empty host": {
			os: func() OS {
				osMock := NewOSMock(mc)
				osMock.GetenvMock.When(pyroscopeHostEnvName).Then("")
				osMock.GetenvMock.When(pyroscopePortEnvName).Then("4040")

				return osMock
			},
			expectedNilObj: false,
			expectedErr:    nil,
		},
		"create config with empty port": {
			os: func() OS {
				osMock := NewOSMock(mc)
				osMock.GetenvMock.When(pyroscopeHostEnvName).Then("127.0.0.1")
				osMock.GetenvMock.When(pyroscopePortEnvName).Then("")

				return osMock
			},
			expectedNilObj: false,
			expectedErr:    nil,
		},
		"create config with empty host and port": {
			os: func() OS {
				osMock := NewOSMock(mc)
				osMock.GetenvMock.When(pyroscopeHostEnvName).Then("")
				osMock.GetenvMock.When(pyroscopePortEnvName).Then("")

				return osMock
			},
			expectedNilObj: false,
			expectedErr:    nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cfg, err := NewPyroscopeConfig(test.os())
			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expectedErr, err)

			if test.expectedNilObj {
				assert.Nil(t, cfg)
			} else {
				assert.NotNil(t, cfg)
			}
		})
	}
}

func TestPyroscopeConfig_Address(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)

	tests := map[string]struct {
		os func() OS

		expectedString string
	}{
		"get address": {
			os: func() OS {
				osMock := NewOSMock(mc)
				osMock.GetenvMock.When(pyroscopeHostEnvName).Then("127.0.0.1")
				osMock.GetenvMock.When(pyroscopePortEnvName).Then("4040")

				return osMock
			},
			expectedString: "127.0.0.1:4040",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cfg, err := NewPyroscopeConfig(test.os())
			assert.Nil(t, err)
			assert.NotNil(t, cfg)
			assert.Equal(t, test.expectedString, cfg.Address())
		})
	}
}
