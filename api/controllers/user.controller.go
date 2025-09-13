package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/patil-prathamesh/yt-backend-go/api/db"
	"github.com/patil-prathamesh/yt-backend-go/api/models"
	"github.com/patil-prathamesh/yt-backend-go/api/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func RegisterUser(c *gin.Context) {
	// get user details from FE``
	// validation - not empty
	// check if user already exists: username and email
	// check for images, check for avatar
	// upload them to cloudinary, avatar
	// create user object - create entry in db
	// remove password and refresh token field from response
	// check for user creation
	// return response

	var user models.User
	collection := db.GetCollection("users")

	username := c.PostForm("username")
	user.Username = strings.ToLower(username)

	user.Email = c.PostForm("email")
	user.FullName = c.PostForm("fullName")
	password := c.PostForm("password")

	if user.Username == "" || user.Email == "" || user.FullName == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "All fields are required",
		})
		return
	}

	filter := bson.M{
		"$or": []bson.M{
			{"email": user.Email},
			{"username": user.Username},
		},
	}

	count, err := db.GetCollection("users").CountDocuments(context.Background(), filter)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if count > 0 {
		c.JSON(409, gin.H{"error": "User with email or username already exists."})
		return
	}

	user.HashPassword(password)

	if file, header, err := c.Request.FormFile("avatar"); err == nil {
		defer file.Close()

		cloudinaryService, err := utils.NewCloudinaryService()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize upload service"})
			return
		}

		filename := fmt.Sprintf("avatar_%s_%s", user.Username, header.Filename)
		avatarURL, err := cloudinaryService.UploadImage(file, filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload avatar"})
			return
		}
		user.Avatar = avatarURL
	}

	if file, header, err := c.Request.FormFile("coverImage"); err == nil {
		defer file.Close()

		cloudinaryService, err := utils.NewCloudinaryService()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize upload service"})
			return
		}

		filename := fmt.Sprintf("cover_%s_%s", user.Username, header.Filename)
		coverURL, err := cloudinaryService.UploadImage(file, filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload cover image"})
			return
		}

		user.CoverImage = coverURL
	}

	user.PrepareForDB()

	_, err = collection.InsertOne(context.Background(), user)

	fmt.Printf("User created with ID: %v\n", user.ID.Hex())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(201, gin.H{
		"message": "User registered successfully",
		"user":    user,
		"success": true,
	})

}
