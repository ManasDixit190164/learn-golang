package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/manasdixit/url-shortener/internal/domain"
	"github.com/manasdixit/url-shortener/internal/middleware"
	"github.com/manasdixit/url-shortener/internal/service"
	"github.com/manasdixit/url-shortener/pkg/response"
)

type URLHandler struct {
	urlService *service.URLService
	logger     *slog.Logger
}

func NewURLHandler(urlService *service.URLService, logger *slog.Logger) *URLHandler {
	return &URLHandler{urlService: urlService, logger: logger}
}

func (h *URLHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req domain.CreateURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	res, err := h.urlService.Create(r.Context(), userID, req)
	if err != nil {
		writeURLError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, response.APIResponse{
		Success: true,
		Message: "short url created successfully",
		Data:    res,
	})
}

func (h *URLHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	res, err := h.urlService.List(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.JSON(w, http.StatusOK, response.APIResponse{Success: true, Data: res})
}

func (h *URLHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := parseUUIDParam(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid url id")
		return
	}

	res, err := h.urlService.GetByID(r.Context(), userID, id)
	if err != nil {
		writeURLError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, response.APIResponse{Success: true, Data: res})
}

func (h *URLHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := parseUUIDParam(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid url id")
		return
	}

	var req domain.UpdateURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	res, err := h.urlService.Update(r.Context(), userID, id, req)
	if err != nil {
		writeURLError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, response.APIResponse{
		Success: true,
		Message: "url updated successfully",
		Data:    res,
	})
}

func (h *URLHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := parseUUIDParam(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid url id")
		return
	}

	if err := h.urlService.Delete(r.Context(), userID, id); err != nil {
		writeURLError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, response.APIResponse{
		Success: true,
		Message: "url deactivated successfully",
	})
}

func (h *URLHandler) Analytics(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := parseUUIDParam(r, "id")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid url id")
		return
	}

	res, err := h.urlService.Analytics(r.Context(), userID, id)
	if err != nil {
		writeURLError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, response.APIResponse{Success: true, Data: res})
}

func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "shortCode")
	if strings.TrimSpace(shortCode) == "" {
		response.Error(w, http.StatusBadRequest, "invalid short code")
		return
	}

	resolvedURL, err := h.urlService.ResolveRedirect(r.Context(), shortCode)
	if err != nil {
		writeURLError(w, err)
		return
	}

	click := domain.Click{
		URLID:     resolvedURL.ID,
		IPAddress: clientIP(r),
		UserAgent: r.UserAgent(),
		Referrer:  r.Referer(),
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := h.urlService.TrackClick(ctx, click); err != nil {
			h.logger.Warn("failed to track click", "error", err, "url_id", resolvedURL.ID)
		}
	}()

	http.Redirect(w, r, resolvedURL.OriginalURL, http.StatusFound)
}

func parseUUIDParam(r *http.Request, key string) (uuid.UUID, error) {
	value := chi.URLParam(r, key)
	return uuid.Parse(value)
}

func clientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		parts := strings.Split(ip, ",")
		return strings.TrimSpace(parts[0])
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func writeURLError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidURL):
		response.Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, service.ErrInvalidInput):
		response.Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, service.ErrAliasUnavailable):
		response.Error(w, http.StatusConflict, err.Error())
	case errors.Is(err, service.ErrURLNotFound):
		response.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, service.ErrURLInactive):
		response.Error(w, http.StatusGone, err.Error())
	case errors.Is(err, service.ErrURLExpired):
		response.Error(w, http.StatusGone, err.Error())
	default:
		response.Error(w, http.StatusInternalServerError, "internal server error")
	}
}
