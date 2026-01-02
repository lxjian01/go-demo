package cmd

import (
	"fmt"
	"go-demo/internal/config"
	"go-demo/internal/logger"
	"go-demo/internal/postgresclient"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	"go-demo/httpd"
)

var rootCmd = &cobra.Command{
	Use:   "go-demo",
	Short: "go-demo http service",
	RunE: func(cmd *cobra.Command, args []string) error {
		conf := config.GetAppConfig()

		// === Gin mode（必须最早）===
		switch conf.Env {
		case "prod":
			gin.SetMode(gin.ReleaseMode)
		case "test":
			gin.SetMode(gin.TestMode)
		default:
			gin.SetMode(gin.DebugMode)
		}

		// === Init Logger ===
		if err := logger.InitLogger(conf.Logger); err != nil {
			return err
		}
		logger.GetLogger().Info().Str("env", conf.Env).Msg("logger initialized")

		// === Init Postgres ===
		if err := postgresclient.InitPostgres(conf.Postgres); err != nil {
			logger.GetLogger().Fatal().Err(err).Msg("postgres init failed")
			return err
		}
		logger.GetLogger().Info().Msg("postgres initialized")

		// === Start HTTP Server（阻塞，内部优雅关闭）===
		httpd.StartHttpServer(conf.Httpd)

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// 初始化配置文件转化成对应的结构体
	cobra.OnInitialize(initConfig)
}
