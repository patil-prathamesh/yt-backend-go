package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/patil-prathamesh/yt-backend-go/api/models"
)

func VerifyJWT(c *gin.Context) {
	var tokenString string
	tokenString, _ = c.Cookie("accss_token")
	if tokenString == "" {
		authHeader := c.GetHeader("Authorization")
		tokenString = strings.Replace(authHeader, "Bearer ", "", 1)
	}

	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		c.Abort()
		return
	}

	secret := os.Getenv("ACCESS_TOKEN_SECRET")
	if secret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server misconfiguration"})
		c.Abort()
		return
	}

	token, err := jwt.ParseWithClaims(tokenString, &models.AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
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

	if claims, ok := token.Claims.(*models.AccessTokenClaims); ok && token.Valid {
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
	}

	c.Next()
}
