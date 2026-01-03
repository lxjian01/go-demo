package logger

import (
	"fmt"
	"go-demo/internal/config"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

var logger zerolog.Logger

func InitLogger(c *config.LoggerConfig) error {
	// 获取日志文件路径
	logFilePath, err := getLogFilePath(c.Dir)
	if err != nil {
		return err
	}

	// 使用 lumberjack 实现日志轮转
	logFile := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    10,   // 最大文件大小 MB
		MaxBackups: 3,    // 最多保留 3 个备份文件
		MaxAge:     28,   // 保留最近 28 天的日志
		Compress:   true, // 压缩旧日志
	}

	// 使用 zerolog 初始化日志记录器
	logger = zerolog.New(logFile).With().Timestamp().Logger()
	// 设置日志等级（例如 Debug）
	parsedLevel, err := zerolog.ParseLevel(c.Level)
	if err != nil {
		return fmt.Errorf("invalid log level: %s", c.Level)
	}
	zerolog.SetGlobalLevel(parsedLevel) // 可以根据需要设置为 DebugLevel 或其他
	return nil
}

// 获取日志文件路径，按日期分割
func getLogFilePath(logDir string) (string, error) {
	date := time.Now().Format("2006-01-02")
	if logDir == "" {
		return "", fmt.Errorf("log dir is empty")
	}
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		return "", fmt.Errorf("mkdir all log dir err %w", err)
	}
	logFile := fmt.Sprintf("%s/app-%s.log", logDir, date)
	return logFile, nil
}

func GetLogger() *zerolog.Logger {
	return &logger
}
