package main

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/patil-prathamesh/yt-backend-go/api/db"
	"github.com/patil-prathamesh/yt-backend-go/api/routes"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err.Error())
	}
	db.ConnectDB()
}

func main() {
	app := gin.New()

	corsConfig := cors.Config{
		AllowOrigins:     []string{os.Getenv("CORS_ORIGIN")},
		AllowCredentials: true,
	}

	app.Use(cors.New(corsConfig))

	routes.SetupRoutes(app)

	app.Run(":" + os.Getenv("PORT"))
}
