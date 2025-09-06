package main

import (
	"github.com/joho/godotenv"
	"github.com/patil-prathamesh/yt-backend-go/api/db"
)

func init() {
	if err := godotenv.Load("../../.env"); err != nil {
		panic(err.Error())
	}
}

func main() {
	db.ConnectDB()
}
