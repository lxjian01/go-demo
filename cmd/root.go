package cmd

import (
	"fmt"
	"go-demo/internal/config"
	"go-demo/internal/logger"
	"go-demo/internal/postgresclient"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"

	"go-demo/httpd"
	"os"
	"path/filepath"
)

// 定义根命令
var rootCmd = &cobra.Command{
	Use: "go-demo",
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.GetAppConfig()
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
				os.Exit(1)
			}
		}()
		fmt.Println("Starting init logger")
		err := logger.InitLogger(conf.Logger)
		if err != nil {
			panic(err)
		}
		fmt.Println("Init logger ok")

		//fmt.Println("Starting init mysql client")
		//err = mysqlclient.InitMysql(conf.Mysql)
		//if err != nil {
		//	panic(err)
		//}
		//fmt.Println("Init mysql client ok")

		fmt.Println("Starting init postgres client")
		err = postgresclient.InitPostgres(conf.Postgres)
		if err != nil {
			panic(err)
		}
		fmt.Println("Init postgres client ok")

		//fmt.Println("Starting init redis client")
		//err = redisclient.InitRedis(conf.Redis)
		//if err != nil {
		//	panic(err)
		//}
		//fmt.Println("Init redis client ok")

		// init gin server
		logger.GetLogger().Info().Msg("Starting init gin server")
		httpd.StartHttpServer(conf.Httpd)
	},
}

func init() {
	// 初始化配置文件转化成对应的结构体
	cobra.OnInitialize(initConfig)
}

// Execute 启动调用的入口方法
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Execute error by ", err)
		os.Exit(1)
	}
}

// initConfig 通过viper初始化配置文件到结构体
func initConfig() {
	dir, _ := os.Getwd()
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}
	configPath := filepath.Join(dir, "configs/"+env)
	// 设置读取的文件路径
	viper.AddConfigPath(configPath)
	// 设置读取的文件名
	viper.SetConfigName("config")
	// 设置文件的类型
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("Read config error by %v \n", err))
	}
	var appConf config.AppConfig
	if err := viper.Unmarshal(&appConf); err != nil {
		panic(fmt.Sprintf("Unmarshal config error by %v \n", err))
	}
	config.SetAppConfig(&appConf)
}
