package handlers

import (
	"encoding/json"
	"net/http"

	"marketplace/internal/auth"
	"marketplace/internal/domain"
	"marketplace/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
	jwtManager  *auth.JWTManager
}

func NewAuthHandler(authService service.AuthService, jwtManager *auth.JWTManager) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		jwtManager:  jwtManager,
	}
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, domain.ErrInvalidInput)
		return
	}

	user, err := h.authService.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	access, _ := h.jwtManager.GenerateAccessToken(user.ID, string(user.Role))
	refresh, _ := h.jwtManager.GenerateRefreshToken(user.ID, string(user.Role))

	writeJSON(w, http.StatusCreated, authResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, domain.ErrInvalidInput)
		return
	}

	user, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	access, _ := h.jwtManager.GenerateAccessToken(user.ID, string(user.Role))
	refresh, _ := h.jwtManager.GenerateRefreshToken(user.ID, string(user.Role))

	writeJSON(w, http.StatusOK, authResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}
