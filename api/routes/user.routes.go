package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/patil-prathamesh/yt-backend-go/api/controllers"
)

func UserRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	// users.Use(middleware.AuthMiddleware())
	{
		users.POST("/register", controllers.RegisterUser)
	}
}
