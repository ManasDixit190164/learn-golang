package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/manasdixit/url-shortener/internal/domain"
	"github.com/manasdixit/url-shortener/internal/service"
	"github.com/manasdixit/url-shortener/pkg/response"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req domain.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	res, err := h.authService.Signup(r.Context(), req)
	if err != nil {
		writeAuthError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, response.APIResponse{
		Success: true,
		Message: "signup successful",
		Data:    res,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	res, err := h.authService.Login(r.Context(), req)
	if err != nil {
		writeAuthError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, response.APIResponse{
		Success: true,
		Message: "login successful",
		Data:    res,
	})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req domain.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	res, err := h.authService.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		writeAuthError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, response.APIResponse{
		Success: true,
		Message: "token refreshed successfully",
		Data:    res,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req domain.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.authService.Logout(r.Context(), req.RefreshToken); err != nil {
		writeAuthError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, response.APIResponse{
		Success: true,
		Message: "logout successful",
	})
}

func writeAuthError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		response.Error(w, http.StatusBadRequest, "name, valid email, and password with minimum 8 characters are required")
	case errors.Is(err, service.ErrEmailAlreadyExists):
		response.Error(w, http.StatusConflict, err.Error())
	case errors.Is(err, service.ErrInvalidCredentials):
		response.Error(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, service.ErrInvalidToken):
		response.Error(w, http.StatusUnauthorized, err.Error())
	default:
		response.Error(w, http.StatusInternalServerError, "internal server error")
	}
}
