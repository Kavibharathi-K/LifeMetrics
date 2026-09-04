package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const tokenExpiration = 24 * time.Hour

func GenerateToken(userID int64) (string, error) {

	secret := os.Getenv("LIFEMETRICS_JWT_SECRET")

	if secret == "" {
		return "", errors.New("LIFEMETRICS_JWT_SECRET is not configured")
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(tokenExpiration).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString([]byte(secret))
}