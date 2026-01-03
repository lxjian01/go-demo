package httpd

import (
	"context"
	"errors"
	"go-demo/internal/redisclient"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"go-demo/httpd/middlewares"
	"go-demo/httpd/routers"
	"go-demo/internal/config"
	"go-demo/internal/logger"
	"go-demo/internal/postgresclient"
)

func StartHttpServer(c *config.HttpdConfig) {
	router := gin.New()
	router.Use(
		gin.Recovery(),
		middlewares.LoggerMiddleware(),
		middlewares.AuthMiddleware(),
	)

	// === K8s probes ===
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/readyz", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		if err := postgresclient.HealthCheck(ctx); err != nil {
			logger.GetLogger().Error().Err(err).Msg("postgres not ready")
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "db not ready"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	routers.UserRoutes(router)

	addr := net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// === 启动 HTTP Server ===
	go func() {
		logger.GetLogger().Info().Str("addr", addr).Msg("http server started")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.GetLogger().Fatal().Err(err).Msg("http server start failed")
		}
	}()

	// === 等待退出信号 ===
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.GetLogger().Info().Str("signal", sig.String()).Msg("shutdown signal received")

	// === 优雅关闭 HTTP ===
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.GetLogger().Error().Err(err).Msg("http server shutdown error")
	} else {
		logger.GetLogger().Info().Msg("http server shutdown complete")
	}

	// === 关闭 PostgreSQL ===
	if err := postgresclient.Close(); err != nil {
		logger.GetLogger().Error().Err(err).Msg("postgres close error")
	} else {
		logger.GetLogger().Info().Msg("postgres closed")
	}

	// === 关闭 PostgreSQL ===
	if err := redisclient.Close(); err != nil {
		logger.GetLogger().Error().Err(err).Msg("redis close error")
	} else {
		logger.GetLogger().Info().Msg("redis closed")
	}
}
