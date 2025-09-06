package db

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/patil-prathamesh/yt-backend-go/api"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectDB() *mongo.Database {
	godotenv.Load("env")
	opts := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
	client, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}
	fmt.Println("Database connected!", client)

	return client.Database(os.Getenv(api.DB_NAME))
}
