package middlewares

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"

	"go-demo/internal/logger"

	"time"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 在日志中记录请求信息
		logger.GetLogger().Info().Str("method", c.Request.Method).Str("url", c.Request.URL.Path).Msg("Request received")
		// 记录请求接收到的时间
		start := time.Now()
		// 捕获请求体（如果需要）
		bodyBytes, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // 重置请求体
		// 记录请求体（可选）
		logger.GetLogger().Debug().Str("request_body", string(bodyBytes)).Msg("request received")
		// 继续处理请求
		c.Next()
		// 计算请求处理耗时
		duration := time.Since(start)
		// 记录响应状态码和耗时
		logger.GetLogger().Info().
			Int("status", c.Writer.Status()).
			Dur("duration", duration).
			Msg("request processed")
	}
}
