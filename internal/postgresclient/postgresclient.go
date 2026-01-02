package postgresclient

import (
	"context"
	"fmt"
	"time"

	"go-demo/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db          *gorm.DB
	sqlDBCloser func() error
)

// InitPostgres 初始化数据库（只允许调用一次）
func InitPostgres(cfg *config.PostgresConfig) error {
	if db != nil {
		return nil
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
		cfg.TimeZone,
	)

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:      nil,
		PrepareStmt: true, // 生产建议开启
	})
	if err != nil {
		return err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return err
	}

	// === 连接池（非常关键）===
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
	sqlDB.SetConnMaxIdleTime(30 * time.Minute)

	// 健康检查
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return err
	}

	db = gormDB
	sqlDBCloser = sqlDB.Close
	return nil
}

// DB 获取数据库实例
func DB() *gorm.DB {
	if db == nil {
		panic("postgres not initialized")
	}
	return db
}

// Close 优雅关闭
func Close() error {
	if sqlDBCloser != nil {
		return sqlDBCloser()
	}
	return nil
}
