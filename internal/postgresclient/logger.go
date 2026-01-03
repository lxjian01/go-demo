package postgresclient

import (
	"context"
	"errors"
	"go-demo/internal/config"
	"time"

	appLogger "go-demo/internal/logger"

	"gorm.io/gorm/logger"
)

type gormZerologger struct {
	logLevel      logger.LogLevel
	slowThreshold time.Duration
}

// newGormLogger 创建 GORM Logger（使用 zerolog）
func newGormLogger(cfg *config.PostgresConfig) logger.Interface {
	level := parseGormLogLevel(cfg.LogLevel)

	return &gormZerologger{
		logLevel:      level,
		slowThreshold: time.Duration(cfg.SlowQueryThresholdMS) * time.Millisecond,
	}
}

func (l *gormZerologger) LogMode(level logger.LogLevel) logger.Interface {
	l.logLevel = level
	return l
}

func (l *gormZerologger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Info {
		appLogger.GetLogger().
			Info().
			Msgf(msg, data...)
	}
}

func (l *gormZerologger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Warn {
		appLogger.GetLogger().
			Warn().
			Msgf(msg, data...)
	}
}

func (l *gormZerologger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Error {
		appLogger.GetLogger().
			Error().
			Msgf(msg, data...)
	}
}

func (l *gormZerologger) Trace(
	ctx context.Context,
	begin time.Time,
	fc func() (sql string, rowsAffected int64),
	err error,
) {
	// 没开启任何日志
	if l.logLevel == logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	log := appLogger.GetLogger().
		With().
		Int64("rows", rows).
		Str("sql", sql).
		Int64("elapsed_ms", elapsed.Milliseconds()).
		Logger()

	switch {
	case err != nil && l.logLevel >= logger.Error && !errors.Is(err, logger.ErrRecordNotFound):
		log.Error().
			Err(err).
			Msg("gorm query error")

	case elapsed > l.slowThreshold && l.slowThreshold != 0 && l.logLevel >= logger.Warn:
		log.Warn().
			Msg("gorm slow query")

	case l.logLevel >= logger.Info:
		log.Info().
			Msg("gorm query")
	}
}

func parseGormLogLevel(level string) logger.LogLevel {
	switch level {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "info":
		return logger.Info
	default:
		return logger.Warn
	}
}
