package logger_test

import (
	"go-demo/internal/config"
	"go-demo/internal/logger"
	"testing"
)

func TestInit(t *testing.T) {
	logConfig := &config.LoggerConfig{
		Dir:   "/Users/lj/logs/go-demo",
		Level: "Info",
	}

	err := logger.InitLogger(logConfig)
	if err != nil {
		t.Errorf("init logger failed: %v", err)
	}
}
