package logger_test

import (
	"go-demo/internal/config"
	"go-demo/internal/logger"
	"testing"
)

func TestInit(t *testing.T) {
	tmpDir := t.TempDir() // 自动创建 & 测试结束自动删除

	logConfig := &config.LoggerConfig{
		Dir:   tmpDir,
		Level: "Info",
	}

	err := logger.InitLogger(logConfig)
	if err != nil {
		t.Fatalf("init logger failed: %v", err)
	}
}
