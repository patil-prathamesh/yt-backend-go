package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/patil-prathamesh/yt-backend-go/api/db"
	"github.com/patil-prathamesh/yt-backend-go/api/routes"
)

func init() {
	if err := godotenv.Load(); err != nil {
        log.Println("Warning: .env file not found")
    }
    db.ConnectDB()
}

func main() {
    // Graceful shutdown
    defer db.DisconnectDB()

    app := gin.New()

    corsConfig := cors.Config{
        AllowOrigins:     []string{os.Getenv("CORS_ORIGIN")},
        AllowCredentials: true,
    }
    app.Use(cors.New(corsConfig))

    routes.SetupRoutes(app)

    // Handle graceful shutdown
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        log.Println("Shutting down gracefully...")
        db.DisconnectDB()
        os.Exit(0)
    }()

    app.Run(":" + os.Getenv("PORT"))
}
