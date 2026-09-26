package health

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Kavibharathi-K/lifemetrics/services/auth"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

func (h *Handler) GoogleConnectHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r.Context())
	if err != nil {

		http.Error(
			w,
			"authenticated user not found",
			http.StatusUnauthorized,
		)
		return
	}
	config := GoogleOAuthConfig()

	// Generate a random state value.
	// This will later be stored against the logged-in user/session
	// and validated in the callback.
	state := uuid.New().String()

	expiresAt := time.Now().Add(10 * time.Minute)

	err = h.Repo.CreateOAuthState(
		r.Context(),
		state,
		userID,
		expiresAt,
	)

	if err != nil {
		http.Error(
			w,
			"failed to create OAuth state",
			http.StatusInternalServerError,
		)
		return
	}

	fmt.Printf(
		"Starting Google Health OAuth for user_id=%d, state=%s\n",
		userID,
		state,
	)

	authURL := config.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("prompt", "consent"),
	)
	fmt.Println("Google OAuth URL:")
	fmt.Println(authURL)

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (h *Handler) GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Check whether Google returned an OAuth error.
	if errParam := r.URL.Query().Get("error"); errParam != "" {
		http.Error(
			w,
			"Google OAuth error: "+errParam,
			http.StatusBadRequest,
		)
		return
	}

	// Get the state returned by Google.
	state := r.URL.Query().Get("state")

	if state == "" {
		http.Error(
			w,
			"Missing OAuth state",
			http.StatusBadRequest,
		)
		return
	}

	// Look up the state and get the LifeMetrics user
	// who started the OAuth flow.
	userID, expiresAt, err := h.Repo.GetOAuthState(
		r.Context(),
		state,
	)

	if err != nil {
		http.Error(
			w,
			"Invalid OAuth state",
			http.StatusBadRequest,
		)
		return
	}

	// Make sure the OAuth state has not expired.
	if time.Now().After(expiresAt) {
		_ = h.Repo.DeleteOAuthState(
			r.Context(),
			state,
		)

		http.Error(
			w,
			"OAuth state expired",
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

	// The OAuth state has successfully completed its purpose.
	// Delete it so it cannot be reused.
	err = h.Repo.DeleteOAuthState(
		r.Context(),
		state,
	)

	if err != nil {
		http.Error(
			w,
			"Failed to consume OAuth state",
			http.StatusInternalServerError,
		)
		return
	}

	// Store the Google OAuth tokens for this LifeMetrics user.
	err = h.Repo.UpsertHealthConnection(
		r.Context(),
		userID,
		"google_health",
		"",
		token.AccessToken,
		token.RefreshToken,
		token.Expiry,
	)

	if err != nil {
		http.Error(
			w,
			"Failed to store Google Health connection",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	fmt.Fprintf(w, "Google Health connected successfully!\n")
}

func (h *Handler) TestTokenHandler(w http.ResponseWriter, r *http.Request) {

	userID, err := auth.GetUserID(r.Context())
	if err != nil {
		http.Error(
			w,
			"authenticated user not found",
			http.StatusUnauthorized,
		)
		return
	}

	accessToken, err := h.GetValidAccessToken(
		r.Context(),
		userID,
	)

	if err != nil {
		http.Error(
			w,
			"failed to get valid access token",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, "Google Health token is valid")
	fmt.Fprintln(w, "Access Token:", accessToken)
}

func (h *Handler) TestStepsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r.Context())
	if err != nil {
		http.Error(
			w,
			"authenticated user not found",
			http.StatusUnauthorized,
		)
		return
	}

	body, err := h.GetSteps(
		r.Context(),
		userID,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
