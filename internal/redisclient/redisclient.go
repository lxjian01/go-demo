package redisclient

import (
	"context"
	"fmt"
	"go-demo/internal/config"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	rdb  *redis.Client
	once sync.Once
)

// InitRedis 初始化 Redis 客户端，返回 error
func InitRedis(cfg *config.RedisConfig) error {
	var initErr error
	once.Do(func() {
		rdb = redis.NewClient(&redis.Options{
			Addr:         cfg.Addr,
			Password:     cfg.Password,
			DB:           cfg.DB,
			PoolSize:     cfg.PoolSize,
			MinIdleConns: cfg.MinIdleConns,
			DialTimeout:  cfg.DialTimeout,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := rdb.Ping(ctx).Err(); err != nil {
			initErr = fmt.Errorf("redis connect failed: %w", err)
			rdb = nil
		}
	})
	return initErr
}

// Close 优雅关闭
func Close() error {
	if rdb != nil {
		return rdb.Close()
	}
	return nil
}

// GetClient 获取全局 Redis 实例
func GetClient() *redis.Client {
	return rdb
}
