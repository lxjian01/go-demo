package config

type AppConfig struct {
	Version string        `yaml:"version"`
	Env     string        `yaml:"env"`
	Logger  *LoggerConfig `yaml:"logger"`
	Httpd   *HttpdConfig  `yaml:"httpd"`
	Mysql   *MysqlConfig  `yaml:"mysql"`
	Redis   *RedisConfig  `yaml:"redisclient"`
}

type HttpdConfig struct {
	Host string
	Port int
}

type LoggerConfig struct {
	Dir   string
	Level string
}

type MysqlConfig struct {
	Host     string
	Port     int
	DbName   string
	User     string
	Password string
	MaxConn  int
	MaxOpen  int
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

var conf *AppConfig

func SetAppConfig(c *AppConfig) {
	conf = c
}

func GetAppConfig() *AppConfig {
	return conf
}
