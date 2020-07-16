package logging

import (
	"github.com/op/go-logging"
	"key-vault-server/config"
)

type Logger struct {
	Logger logging.Logger
	cfg    *config.Config
}

func NewLogger(cfg *config.Config) *Logger {
	logger := logging.MustGetLogger("key-vault-server")
	return &Logger{
		Logger: *logger,
		cfg:    cfg,
	}
}
