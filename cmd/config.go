package cmd

import (
	"go-demo/internal/config"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func initConfig() {
	dir, _ := os.Getwd()

	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}

	configPath := filepath.Join(dir, "configs", env)

	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	var appConf config.AppConfig
	if err := viper.Unmarshal(&appConf); err != nil {
		panic(err)
	}

	// === env 只用于加载配置，行为以配置为准 ===
	appConf.Env = env
	config.SetAppConfig(&appConf)
}
