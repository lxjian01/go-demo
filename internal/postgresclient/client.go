package postgresclient

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go-demo/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db      *gorm.DB
	once    sync.Once
	initErr error

	sqlDBCloser func() error
)

// InitPostgres 初始化数据库（进程级只会执行一次）
func InitPostgres(cfg *config.PostgresConfig) error {
	once.Do(func() {
		initErr = initPostgres(cfg)
	})
	return initErr
}

func initPostgres(cfg *config.PostgresConfig) error {
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
		Logger:      newGormLogger(cfg),
		PrepareStmt: cfg.PrepareStmt, // 由配置控制
	})
	if err != nil {
		return err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return err
	}

	// === 连接池配置（生产关键）===
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.ConnMaxIdleTime) * time.Second)

	// === 启动时健康检查（fail-fast）===
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return err
	}

	db = gormDB
	sqlDBCloser = sqlDB.Close
	return nil
}

// DB 返回 *gorm.DB（禁止在 Init 之前调用）
func DB() *gorm.DB {
	if db == nil {
		panic("postgres client: DB not initialized")
	}
	return db
}

// DBWithCtx 返回带 context 的 DB（Gin / HTTP 必用）
func DBWithCtx(ctx context.Context) *gorm.DB {
	return DB().WithContext(ctx)
}

// Close 优雅关闭数据库连接池
func Close() error {
	if sqlDBCloser != nil {
		return sqlDBCloser()
	}
	return nil
}
