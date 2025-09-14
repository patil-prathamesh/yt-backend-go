package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/patil-prathamesh/yt-backend-go/api/db"
	"github.com/patil-prathamesh/yt-backend-go/api/models"
	"github.com/patil-prathamesh/yt-backend-go/api/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func LoginUser(c *gin.Context) {
	// req body -> data
	// username or email
	// find the user
	// password check
	// access and refresh token
	// send cookie

	type Request struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var request Request
	collection := db.GetCollection("users")

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	filter := bson.M{
		"$or": []bson.M{
			{"email": request.Email},
			{"username": request.Username},
		},
	}

	var user models.User
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	err = user.IsPasswordCorrect(request.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	accessToken, err := user.GenerateAccessToken()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})

		fmt.Println("Access token error:", err.Error())
		return
	}

	refreshToken, err := user.GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})

		fmt.Println("Refresh token error:", err.Error())
		return
	}

	collection.UpdateByID(context.Background(), user.ID, bson.M{
		"$set": bson.M{
			"refreshToken": refreshToken,
			"updatedAt":    time.Now(),
		},
	})

	c.SetCookie("access_token", accessToken, 3600*1, "/", "", true, true)
	c.SetCookie("refresh_token", refreshToken, 3600*24*10, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"id":          user.ID.Hex(),
			"username":    user.Username,
			"email":       user.Email,
			"accessToken": accessToken,
		},
	})
}

func LogoutUser(c *gin.Context) {
	id, _ := c.Get("user_id")
	objectId, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%v", id))
	fmt.Println(id)
	collection := db.GetCollection("users")
	result, err := collection.UpdateByID(context.Background(), objectId, bson.M{
		"$set": bson.M{
			"refreshToken": "",
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while logout"})
		fmt.Println(err.Error())
		return
	}

	c.SetCookie("access_token", "", -1, "/", "", true, true)
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"message":        "Logout successful",
		"modified_count": result.ModifiedCount,
	})
}

func RefreshAccessToken(c *gin.Context) {
	var tokenString string
	tokenString, _ = c.Cookie("refresh_token")
	if tokenString == "" {
		authHeader := c.GetHeader("Authorization")
		tokenString = strings.Replace(authHeader, "Bearer ", "", 1)
	}

	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	secret := os.Getenv("REFRESH_TOKEN_SECRET")
	if secret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server misconfiguration"})
		return
	}

	token, err := jwt.ParseWithClaims(tokenString, &models.RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		c.Abort()
		return
	}

	var id string

	if claims, ok := token.Claims.(*models.RefreshTokenClaims); ok && token.Valid {
		id = claims.UserID
	}
	objectId, _ := primitive.ObjectIDFromHex(fmt.Sprintf("%v", id))

	var user models.User
	collection := db.GetCollection("users")

	collection.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&user)

	var accessToken string

	if tokenString == user.RefreshToken {
		accessToken, _ = user.GenerateAccessToken()
		c.SetCookie("access_token", accessToken, 3600, "/", "", true, true)
	}

	c.JSON(200, gin.H{"access_token": accessToken})
}

func ChangeCurrentPassword(c *gin.Context) {
	type Request struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}

	var request Request
	collection := db.GetCollection("users")

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	userId, _ := c.Get("user_id")
	userObjectId := fmt.Sprintf("%v", userId)

	var user models.User
	collection.FindOne(context.Background(), bson.M{"_id": userObjectId}).Decode(&user)

	err := user.IsPasswordCorrect(request.OldPassword)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect old password"})
		return
	}

	result, err := collection.UpdateByID(context.Background(), userObjectId, bson.M{
		"$set": bson.M{
			"password": request.NewPassword,
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while updating password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password updated successfully", "count": result.ModifiedCount})
}

func GetCurrentUser(c *gin.Context) {
	collection := db.GetCollection("users")
	userName, _ := c.Get("username")
	var user models.User

	collection.FindOne(context.Background(), bson.M{"username": userName}).Decode(&user)

	c.JSON(http.StatusOK, gin.H{"message": "Current user fetched succssfully", "data": user})
}

func UpdateUserAvatar(c *gin.Context) {
	userName, _ := c.Get("username")
	collection := db.GetCollection("users")
	if file, header, err := c.Request.FormFile("avatar"); err == nil {
		defer file.Close()

		cloudinaryService, err := utils.NewCloudinaryService()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize upload service"})
			return
		}

		filename := fmt.Sprintf("avatar_%s_%s", userName, header.Filename)
		avatarURL, err := cloudinaryService.UploadImage(file, filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload avatar"})
			return
		}

		if avatarURL == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating avatar"})
			return
		}

		result, err := collection.UpdateOne(context.Background(), bson.M{"username": userName}, bson.M{
			"$set": bson.M{
				"avatar": avatarURL,
			},
		})

		if err != nil || result.MatchedCount == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating avatar"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Avatar updated successfully"})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "Avatar is missing"})
}

func UpdateUserCoverImage(c *gin.Context) {
	userName, _ := c.Get("username")
	collection := db.GetCollection("users")
	if file, header, err := c.Request.FormFile("avatar"); err == nil {
		defer file.Close()

		cloudinaryService, err := utils.NewCloudinaryService()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize upload service"})
			return
		}

		filename := fmt.Sprintf("cover_%s_%s", userName, header.Filename)
		coverURL, err := cloudinaryService.UploadImage(file, filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload cover image"})
			return
		}

		if coverURL == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating cover image"})
			return
		}

		result, err := collection.UpdateOne(context.Background(), bson.M{"username": userName}, bson.M{
			"$set": bson.M{
				"coverImage": coverURL,
			},
		})

		if err != nil || result.MatchedCount == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating cover image"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Cover Image updated successfully"})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "Cover image is missing"})
}
