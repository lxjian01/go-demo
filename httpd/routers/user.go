package routers

import (
	"go-demo/httpd/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(route *gin.Engine) {
	user := route.Group("/api/user")
	{
		// alarm group
		user.POST("/", controllers.AddUser)
		user.GET("/page", controllers.GetUserPage)
	}
}
