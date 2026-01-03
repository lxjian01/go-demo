package config

import "time"

type AppConfig struct {
	Version  string          `mapstructure:"version"`
	Env      string          `mapstructure:"env"`
	Logger   *LoggerConfig   `mapstructure:"logger"`
	Httpd    *HttpdConfig    `mapstructure:"httpd"`
	Postgres *PostgresConfig `mapstructure:"postgres"`
	Redis    *RedisConfig    `mapstructure:"redis"`
}

type HttpdConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type LoggerConfig struct {
	Dir   string `mapstructure:"dir"`
	Level string `mapstructure:"level"`
}

type PostgresConfig struct {
	Host                 string        `mapstructure:"host"`                 // 数据库地址
	Port                 int           `mapstructure:"port"`                 // 端口
	User                 string        `mapstructure:"user"`                 // 用户名
	Password             string        `mapstructure:"password"`             // 密码
	DB                   string        `mapstructure:"db"`                   // 数据库名
	SSLMode              string        `mapstructure:"ssl_mode"`             // ssl 模式，内网可用 require，公网建议 verify-full
	Timezone             string        `mapstructure:"timezone"`             // 时区
	MaxIdleConns         int           `mapstructure:"max_idle_conns"`       // 最大空闲连接数
	MaxOpenConns         int           `mapstructure:"max_open_conns"`       // 最大打开连接数
	ConnMaxLifetime      time.Duration `mapstructure:"conn_max_lifetime"`    // 连接最大生命周期
	ConnMaxIdleTime      time.Duration `mapstructure:"conn_max_idle_time"`   // 连接最大空闲时间
	PrepareStmt          bool          `mapstructure:"prepare_stmt"`         // 是否启用预编译
	LogLevel             string        `mapstructure:"logLevel"`             // GORM 日志级别
	SlowQueryThresholdMS int           `mapstructure:"slow_query_threshold"` // 慢 SQL 阈值，ms
}

type RedisConfig struct {
	Addr         string        `mapstructure:"addr"`           // Redis 地址
	Password     string        `mapstructure:"password"`       // Redis 密码，如果为空则无密码
	DB           int           `mapstructure:"db"`             // Redis 数据库索引，默认 0
	PoolSize     int           `mapstructure:"pool_size"`      // 连接池最大连接数
	MinIdleConns int           `mapstructure:"min_idle_conns"` // 连接池最小空闲连接数，避免频繁创建连接
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`   // 连接超时时间，建议 3~5s
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`   // 读取超时时间，建议 2~3s
	WriteTimeout time.Duration `mapstructure:"write_timeout"`  // 写入超时时间，建议 2~3s
}

var conf *AppConfig

func SetAppConfig(c *AppConfig) {
	conf = c
}

func GetAppConfig() *AppConfig {
	return conf
}
