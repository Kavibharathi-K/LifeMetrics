package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	Repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{
		Repo: repo,
	}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {

	Email    string `json:"email"`
	Password string `json:"password"`

}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {

	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if req.Email == "" || req.Password == "" {
		http.Error(
			w,
			"email and password are required",
			http.StatusBadRequest,
		)
		return
	}

	if len(req.Password) < 8 {
		http.Error(
			w,
			"password must be at least 8 characters",
			http.StatusBadRequest,
		)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		http.Error(
			w,
			"failed to process password",
			http.StatusInternalServerError,
		)
		return
	}

	userID, err := h.Repo.CreateUser(
		r.Context(),
		req.Email,
		string(passwordHash),
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

	json.NewEncoder(w).Encode(
		map[string]interface{}{
			"user_id": userID,
			"email":   req.Email,
		},
	)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

	var req LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if req.Email == "" || req.Password == "" {
		http.Error(
			w,
			"email and password are required",
			http.StatusBadRequest,
		)
		return
	}

	user, err := h.Repo.GetUserByEmail(
		r.Context(),
		req.Email,
	)

	if err != nil {
		http.Error(
			w,
			"invalid email or password",
			http.StatusUnauthorized,
		)
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	)

	if err != nil {
		http.Error(
			w,
			"invalid email or password",
			http.StatusUnauthorized,
		)
		return
	}

	token, err := GenerateToken(user.ID)

	if err != nil {
		http.Error(
			w,
			"failed to generate authentication token",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(
		map[string]interface{}{
			"user_id": user.ID,
			"email":   user.Email,
			"token":   token,
		},
	)
}