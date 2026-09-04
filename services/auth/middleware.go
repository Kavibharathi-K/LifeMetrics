package auth

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userIDKey contextKey = "user_id"

func AuthMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(
				w,
				"authorization header required",
				http.StatusUnauthorized,
			)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(
				w,
				"invalid authorization header",
				http.StatusUnauthorized,
			)
			return
		}

		tokenString := parts[1]

		secret := os.Getenv("LIFEMETRICS_JWT_SECRET")

		if secret == "" {
			http.Error(
				w,
				"authentication configuration error",
				http.StatusInternalServerError,
			)
			return
		}

		token, err := jwt.Parse(
			tokenString,
			func(token *jwt.Token) (interface{}, error) {

				if token.Method != jwt.SigningMethodHS256 {
					return nil, errors.New("unexpected signing method")
				}

				return []byte(secret), nil
			},
			jwt.WithValidMethods([]string{
				jwt.SigningMethodHS256.Alg(),
			}),
		)

		if err != nil || !token.Valid {
			http.Error(
				w,
				"invalid or expired token",
				http.StatusUnauthorized,
			)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			http.Error(
				w,
				"invalid token claims",
				http.StatusUnauthorized,
			)
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)

		if !ok {
			http.Error(
				w,
				"invalid user id in token",
				http.StatusUnauthorized,
			)
			return
		}

		userID := int64(userIDFloat)

		ctx := context.WithValue(
			r.Context(),
			userIDKey,
			userID,
		)

		next.ServeHTTP(
			w,
			r.WithContext(ctx),
		)
	})
}

func GetUserID(ctx context.Context) (int64, error) {

	userID, ok := ctx.Value(userIDKey).(int64)

	if !ok {
		return 0, errors.New("user id not found in context")
	}

	return userID, nil
}
