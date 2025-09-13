package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/patil-prathamesh/yt-backend-go/api"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Client *mongo.Client
var Database *mongo.Database

func ConnectDB() {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MONGODB_URI environment variable not set")
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	opts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(opts)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Test the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	fmt.Println("✅ MongoDB connected successfully!")

	// Set global variables
	Client = client
	Database = client.Database(api.DB_NAME)
}

func GetCollection(collectionName string) *mongo.Collection {
	return Database.Collection(collectionName)
}

func DisconnectDB() {
	if Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := Client.Disconnect(ctx); err != nil {
			log.Println("Error disconnecting from MongoDB:", err)
		} else {
			fmt.Println("✅ MongoDB disconnected successfully!")
		}
	}
}
