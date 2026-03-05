package handlers

import (
	"encoding/json"
	"net/http"
	"smart-todo/internal/api/dto"
	"smart-todo/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	authResp, err := h.authService.Register(r.Context(), &service.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		switch err {
		case service.ErrUserExists:
			http.Error(w, "user already exists", http.StatusConflict)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	resp := dto.AuthResponse{
		Token: authResp.Token,
		User: &dto.UserResponse{
			ID:        authResp.User.ID,
			Username:  authResp.User.Username,
			Email:     authResp.User.Email,
			CreatedAt: authResp.User.CreatedAt.Format("2006-01-02T15:04:05Z"),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	authResp, err := h.authService.Login(r.Context(), &service.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		switch err {
		case service.ErrInvalidCredentials:
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	resp := dto.AuthResponse{
		Token: authResp.Token,
		User: &dto.UserResponse{
			ID:        authResp.User.ID,
			Username:  authResp.User.Username,
			Email:     authResp.User.Email,
			CreatedAt: authResp.User.CreatedAt.Format("2006-01-02T15:04:05Z"),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
