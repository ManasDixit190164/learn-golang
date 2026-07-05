package domain

import (
	"time"

	"github.com/google/uuid"
)

type ShortURL struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	OriginalURL string     `json:"original_url"`
	ShortCode   string     `json:"short_code"`
	CustomAlias *string    `json:"custom_alias,omitempty"`
	Title       *string    `json:"title,omitempty"`
	IsActive    bool       `json:"is_active"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateURLRequest struct {
	OriginalURL string     `json:"original_url"`
	CustomAlias *string    `json:"custom_alias"`
	Title       *string    `json:"title"`
	ExpiresAt   *time.Time `json:"expires_at"`
}

type UpdateURLRequest struct {
	OriginalURL *string    `json:"original_url"`
	Title       *string    `json:"title"`
	IsActive    *bool      `json:"is_active"`
	ExpiresAt   *time.Time `json:"expires_at"`
}

type URLResponse struct {
	ID          uuid.UUID  `json:"id"`
	OriginalURL string     `json:"original_url"`
	ShortCode   string     `json:"short_code"`
	ShortURL    string     `json:"short_url"`
	CustomAlias *string    `json:"custom_alias,omitempty"`
	Title       *string    `json:"title,omitempty"`
	IsActive    bool       `json:"is_active"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func NewURLResponse(url ShortURL, baseURL string) URLResponse {
	return URLResponse{
		ID:          url.ID,
		OriginalURL: url.OriginalURL,
		ShortCode:   url.ShortCode,
		ShortURL:    baseURL + "/" + url.ShortCode,
		CustomAlias: url.CustomAlias,
		Title:       url.Title,
		IsActive:    url.IsActive,
		ExpiresAt:   url.ExpiresAt,
		CreatedAt:   url.CreatedAt,
		UpdatedAt:   url.UpdatedAt,
	}
}
