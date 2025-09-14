package models

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	WatchHistory []primitive.ObjectID `json:"watchHistory" bson:"watchHistory"`
	Username     string               `json:"username" bson:"username"`
	Email        string               `json:"email" bson:"email"`
	FullName     string               `json:"fullName" bson:"fullName"`
	Avatar       string               `json:"avatar" bson:"avatar"`
	CoverImage   string               `json:"coverImage" bson:"coverImage"`
	Password     string               `json:"-" bson:"password"`
	RefreshToken string               `json:"-" bson:"refreshToken"`
	CreatedAt    time.Time            `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time            `json:"updatedAt" bson:"updatedAt"`
}

// JWT Claims structure
type AccessTokenClaims struct {
	UserID   string `json:"_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	FullName string `json:"fullName"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	UserID string `json:"_id"`
	jwt.RegisteredClaims
}

func (u *User) PrepareForDB() {
    if u.ID.IsZero() {
        u.ID = primitive.NewObjectID()
    }
    
    now := time.Now()
    
    // Set CreatedAt only for new records
    if u.CreatedAt.IsZero() {
        u.CreatedAt = now
    }
    
    // Always update UpdatedAt
    u.UpdatedAt = now
}

func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) IsPasswordCorrect(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (u *User) GenerateAccessToken() (string, error) {
	// Parse expiry from environment variable (e.g., "1d", "24h")
	expiryStr := os.Getenv("ACCESS_TOKEN_EXPIRY")
	if expiryStr == "" {
		expiryStr = "15m" // Default 15 minutes
	}

	// Convert expiry string to duration
	var expiry time.Duration
	if expiryStr == "1d" {
		expiry = 24 * time.Hour
	} else {
		var err error
		expiry, err = time.ParseDuration(expiryStr)
		if err != nil {
			expiry = 15 * time.Minute // Fallback to 15 minutes
		}
	}

	// Create claims
	claims := AccessTokenClaims{
		UserID:   u.ID.Hex(),
		Email:    u.Email,
		Username: u.Username,
		FullName: u.FullName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "yt-backend",
			Subject:   u.ID.Hex(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Get secret from environment
	secret := os.Getenv("ACCESS_TOKEN_SECRET")
	if secret == "" {
		return "", jwt.ErrInvalidKey
	}

	// Sign and return token string
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (u *User) GenerateRefreshToken() (string, error) {
	// Parse expiry from environment variable
	expiryStr := os.Getenv("REFRESH_TOKEN_EXPIRY")
	if expiryStr == "" {
		expiryStr = "7d" // Default 7 days
	}

	// Convert expiry string to duration
	var expiry time.Duration
	if expiryStr == "10d" {
		expiry = 10 * 24 * time.Hour
	} else {
		var err error
		expiry, err = time.ParseDuration(expiryStr)
		if err != nil {
			expiry = 7 * 24 * time.Hour // Fallback to 7 days
		}
	}

	// Create claims (minimal for refresh token)
	claims := RefreshTokenClaims{
		UserID: u.ID.Hex(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "yt-backend",
			Subject:   u.ID.Hex(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Get secret from environment
	secret := os.Getenv("REFRESH_TOKEN_SECRET")
	if secret == "" {
		return "", jwt.ErrInvalidKey
	}

	// Sign and return token string
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
