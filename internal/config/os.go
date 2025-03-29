package config

import "os"

// OS ...
type OS interface {
	GetEnv(key string) string
}

type osImpl struct{}

// NewOS ...
func NewOS() OS {
	return &osImpl{}
}

// GetEnv ...
func (osImpl) GetEnv(key string) string {
	return os.Getenv(key)
}
