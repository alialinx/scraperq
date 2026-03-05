package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alialin/scraperq/internal/auth"
	"github.com/alialin/scraperq/internal/models"
	"github.com/alialin/scraperq/internal/repository"
	"github.com/google/uuid"
)

type AuthHandler struct {
	repo      *repository.UserRepo
	jwtSecret string
}

func NewAuthHandler(repo *repository.UserRepo, jwtSecret string) *AuthHandler {
	return &AuthHandler{repo: repo, jwtSecret: jwtSecret}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

	var req models.RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	passwordHash, err := auth.HashPassword(req.Password)

	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	apiKey := uuid.New().String()

	user := models.User{
		Email:        req.Email,
		PasswordHash: passwordHash,
		APIKey:       apiKey,
		DailyLimit:   100,
		MonthlyLimit: 2000,
	}

	err = h.repo.Create(r.Context(), &user)
	if err == repository.ErrEmailExists {
		http.Error(w, "email already exists", http.StatusConflict)
		return
	}

	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Email, h.jwtSecret)

	if err != nil {
		http.Error(w, "token error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.AuthResponse{
		Token:  token,
		APIKey: user.APIKey,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	var req models.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.repo.FindByEmail(r.Context(), req.Email)

	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	isTrue := auth.CheckPassword(req.Password, user.PasswordHash)

	if !isTrue {
		http.Error(w, "wrong password", http.StatusUnauthorized)
		return
	}
	token, err := auth.GenerateToken(user.ID, user.Email, h.jwtSecret)
	if err != nil {
		http.Error(w, "token error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.AuthResponse{
		Token:  token,
		APIKey: user.APIKey,
	})

}
