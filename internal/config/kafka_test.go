package config

import (
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestNewKafkaConfig(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)

	tests := map[string]struct {
		os func() OS

		expectedNilObj bool
		expectedErr    error
		expectedConfig *KafkaConfig
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
				osMock.GetenvMock.When(kafkaHostEnvName + "1").Then("127.0.0.1:9092")
				osMock.GetenvMock.When(kafkaHostEnvName + "2").Then("127.0.0.1:9093")
				osMock.GetenvMock.When(kafkaHostEnvName + "3").Then("127.0.0.1:9094")
				osMock.GetenvMock.When(kafkaHostEnvName + "4").Then("")

				return osMock
			},

			expectedNilObj: false,
			expectedErr:    nil,
			expectedConfig: &KafkaConfig{
				addresses: []string{"127.0.0.1:9092", "127.0.0.1:9093", "127.0.0.1:9094"},
				enabled:   true,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cfg, err := NewKafkaConfig(test.os())
			assert.Equal(t, test.expectedErr, err)

			if test.expectedNilObj {
				assert.Nil(t, cfg)
			} else {
				assert.NotNil(t, cfg)
				assert.Equal(t, test.expectedConfig, cfg)
			}
		})
	}
}

func TestGetAddresses(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)

	tests := map[string]struct {
		os func() OS

		expectedAddresses []string
	}{
		"get empty addresses": {
			os: func() OS {
				osMock := NewOSMock(mc)
				osMock.GetenvMock.When(kafkaHostEnvName + "1").Then("")

				return osMock
			},

			expectedAddresses: []string{},
		},
		"get addresses": {
			os: func() OS {
				osMock := NewOSMock(mc)
				osMock.GetenvMock.When(kafkaHostEnvName + "1").Then("127.0.0.1:9092")
				osMock.GetenvMock.When(kafkaHostEnvName + "2").Then("127.0.0.1:9093")
				osMock.GetenvMock.When(kafkaHostEnvName + "3").Then("127.0.0.1:9094")
				osMock.GetenvMock.When(kafkaHostEnvName + "4").Then("")

				return osMock
			},

			expectedAddresses: []string{"127.0.0.1:9092", "127.0.0.1:9093", "127.0.0.1:9094"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cfg, err := NewKafkaConfig(test.os())
			assert.Nil(t, err)
			assert.NotNil(t, cfg)
			assert.Equal(t, test.expectedAddresses, cfg.Addresses())
		})
	}
}
