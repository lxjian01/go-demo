package config

type AppConfig struct {
	Version  string          `yaml:"version"`
	Env      string          `yaml:"env"`
	Logger   *LoggerConfig   `yaml:"logger"`
	Httpd    *HttpdConfig    `yaml:"httpd"`
	Mysql    *MysqlConfig    `yaml:"mysql"`
	Postgres *PostgresConfig `yaml:"postgres"`
	Redis    *RedisConfig    `yaml:"redisclient"`
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

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string // disable / require / verify-full
	TimeZone string
	// === 连接池 ===
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int // seconds
	ConnMaxIdleTime int // seconds
	// === GORM ===
	PrepareStmt        bool
	LogLevel           string // silent / error / warn / info
	SlowQueryThreshold int    // ms
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
