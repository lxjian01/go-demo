package httpd

import (
	"go-demo/httpd/routers"
	"net"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-demo/internal/logger"

	"go-demo/httpd/middlewares"
	"go-demo/internal/config"
)

func StartHttpServer(c *config.HttpdConfig) {
	router := gin.Default()
	// 添加自定义的 logger 间件
	router.Use(middlewares.LoggerMiddleware(), gin.Recovery())
	router.Use(middlewares.AuthMiddleware(), gin.Recovery())
	// Define a simple GET endpoint
	router.GET("/ping", func(c *gin.Context) {
		// Return JSON response
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
		logger.GetLogger().Info().Str("Start HTTP server at %s", "111")
	})
	// 添加路由
	routers.UserRoutes(router) //Added alarm routers
	// 拼接host
	Host := c.Host
	Port := strconv.Itoa(c.Port)
	addr := net.JoinHostPort(Host, Port)
	logger.GetLogger().Info().Str("server", addr).Msg("Start HTTP server")
	err := router.Run(addr)
	if err != nil {
		logger.GetLogger().Error().AnErr("Start HTTP server error by %v", err)
	}
}
