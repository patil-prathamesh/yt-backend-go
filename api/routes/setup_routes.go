package routes

import "github.com/gin-gonic/gin"

func SetupRoutes(app *gin.Engine) {
	v1 := app.Group("/api/v1")

	UserRoutes(v1)
}	