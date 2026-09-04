package health

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

func GoogleConnectHandler(w http.ResponseWriter, r *http.Request) {
	config := GoogleOAuthConfig()

	// Generate a random state value.
	// This will later be stored against the logged-in user/session
	// and validated in the callback.
	state := uuid.New().String()

	authURL := config.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
	)

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Check whether Google returned an OAuth error.
	if errParam := r.URL.Query().Get("error"); errParam != "" {
		http.Error(
			w,
			"Google OAuth error: "+errParam,
			http.StatusBadRequest,
		)
		return
	}

	// Get the authorization code returned by Google.
	code := r.URL.Query().Get("code")

	if code == "" {
		http.Error(
			w,
			"Missing authorization code",
			http.StatusBadRequest,
		)
		return
	}

	config := GoogleOAuthConfig()

	// Exchange the temporary authorization code
	// for an access token and refresh token.
	token, err := config.Exchange(r.Context(), code)
	if err != nil {
		http.Error(
			w,
			"Failed to exchange authorization code: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	// Temporary debugging output.
	// We will NOT keep tokens like this in production.
	w.Header().Set("Content-Type", "text/plain")

	fmt.Fprintf(w, "Google OAuth successful!\n\n")
	fmt.Fprintf(w, "Access Token: %s\n\n", token.AccessToken)
	fmt.Fprintf(w, "Token Type: %s\n\n", token.TokenType)
	fmt.Fprintf(w, "Expiry: %s\n\n", token.Expiry)

	if token.RefreshToken != "" {
		fmt.Fprintf(w, "Refresh Token: %s\n", token.RefreshToken)
	} else {
		fmt.Fprintf(w, "Refresh Token: <not returned>\n")
	}
}
