package redisclient

import (
	"context"
	"fmt"
	"go-demo/internal/config"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	Rdb  *redis.Client
	once sync.Once
)

// InitRedis 初始化 Redis 客户端
func InitRedis(conf *config.RedisConfig) error {
	var err error
	once.Do(func() {
		// 创建 Redis 客户端
		Rdb = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port), // Redis 地址
			Password: conf.Password,                              // Redis 密码
			DB:       conf.DB,                                    // Redis 默认数据库
		})

		// 测试 Redis 连接
		_, err = Rdb.Ping(context.Background()).Result()
		if err != nil {
			return
		}
	})
	return err
}
