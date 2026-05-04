package jwtToken

import (
	"auth/internal/models"
	"fmt"
	"os"
	"time"
)

import "github.com/golang-jwt/jwt"

func New(user models.User, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.UserID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", fmt.Errorf("jwt.New: %w", err)
	}

	return tokenString, nil
}

func Verify(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
}
