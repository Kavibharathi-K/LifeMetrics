package health

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/oauth2"
)

type HealthConnection struct {
	ProviderUserID string
	AccessToken    string
	RefreshToken   string
	ExpiresAt      time.Time
}

func (h *Handler) GetValidAccessToken(
	ctx context.Context,
	userID int64,
) (string, error) {

	connection, err := h.Repo.GetHealthConnection(
		ctx,
		userID,
		"google_health",
	)

	if err != nil {
		return "", err
	}

	// Give the access token a small safety window.
	// This prevents us from using a token that is about to expire.
	if time.Now().Before(connection.ExpiresAt.Add(-2 * time.Minute)) {
		fmt.Println("Token valid and using the same token")
		return connection.AccessToken, nil
	}

	if connection.RefreshToken == "" {
		return "", errors.New("google health refresh token is missing")
	}

	config := GoogleOAuthConfig()

	token := &oauth2.Token{
		AccessToken:  connection.AccessToken,
		RefreshToken: connection.RefreshToken,
		Expiry:       connection.ExpiresAt,
	}

	fmt.Println("Token expired, refreshing the token")
	tokenSource := config.TokenSource(ctx, token)

	newToken, err := tokenSource.Token()
	if err != nil {
		return "", err
	}
	fmt.Println("NEW TOKEN OBJ === ", newToken)

	// Google may not return a new refresh token when refreshing.
	// Preserve the existing one in that case.
	refreshToken := newToken.RefreshToken

	if refreshToken == "" {
		refreshToken = connection.RefreshToken
	}

	err = h.Repo.UpsertHealthConnection(
		ctx,
		userID,
		"google_health",
		connection.ProviderUserID,
		newToken.AccessToken,
		refreshToken,
		newToken.Expiry,
	)

	if err != nil {
		return "", err
	}

	return newToken.AccessToken, nil
}
