package middlewares

import (
	"go-demo/internal/logger"

	"github.com/gin-gonic/gin"
)

type LoginUser struct {
	UserName string `json:"user_name"`
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		loginUser := &LoginUser{
			UserName: c.Request.URL.Path,
		}
		c.Set("loginUser", loginUser)
		// 继续执行请求
		c.Next()
		logger.GetLogger().Info()
	}
}
