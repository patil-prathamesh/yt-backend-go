package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/patil-prathamesh/yt-backend-go/api/controllers"
	"github.com/patil-prathamesh/yt-backend-go/api/middlewares"
)

func UserRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	// users.Use(middleware.AuthMiddleware())
	{
		users.POST("/register", controllers.RegisterUser)
		users.POST("/login", controllers.LoginUser)

		// users.Use(middlewares.VerifyJWT)
		users.POST("/logout", middlewares.VerifyJWT, controllers.LogoutUser)
		users.POST("/refresh-token", controllers.RefreshAccessToken)
		users.POST("/change-password", middlewares.VerifyJWT, controllers.ChangeCurrentPassword)
		users.GET("/current-user", middlewares.VerifyJWT, controllers.GetCurrentUser)
		users.PATCH("/avatar", middlewares.VerifyJWT, controllers.UpdateUserAvatar)
		users.PATCH("/cover-image", middlewares.VerifyJWT, controllers.UpdateUserCoverImage)
		users.GET("/c/:username", middlewares.VerifyJWT, controllers.GetUserChannelProfile)
		users.GET("/history", middlewares.VerifyJWT, controllers.GetWatchHistory)
	}
}
