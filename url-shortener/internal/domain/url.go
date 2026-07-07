package domain // package declaration for the module

import ( // start import block
	"time" // import package

	"github.com/google/uuid" // import package
) // end import block or block scope

type ShortURL struct { // declare struct type
	ID          uuid.UUID  `json:"id"` // execute statement
	UserID      uuid.UUID  `json:"user_id"` // execute statement
	OriginalURL string     `json:"original_url"` // execute statement
	ShortCode   string     `json:"short_code"` // execute statement
	CustomAlias *string    `json:"custom_alias,omitempty"` // execute statement
	Title       *string    `json:"title,omitempty"` // execute statement
	IsActive    bool       `json:"is_active"` // execute statement
	ExpiresAt   *time.Time `json:"expires_at,omitempty"` // execute statement
	CreatedAt   time.Time  `json:"created_at"` // execute statement
	UpdatedAt   time.Time  `json:"updated_at"` // execute statement
} // end block

type CreateURLRequest struct { // declare struct type
	OriginalURL string     `json:"original_url"` // execute statement
	CustomAlias *string    `json:"custom_alias"` // execute statement
	Title       *string    `json:"title"` // execute statement
	ExpiresAt   *time.Time `json:"expires_at"` // execute statement
} // end block

type UpdateURLRequest struct { // declare struct type
	OriginalURL *string    `json:"original_url"` // execute statement
	Title       *string    `json:"title"` // execute statement
	IsActive    *bool      `json:"is_active"` // execute statement
	ExpiresAt   *time.Time `json:"expires_at"` // execute statement
} // end block

type URLResponse struct { // declare struct type
	ID          uuid.UUID  `json:"id"` // execute statement
	OriginalURL string     `json:"original_url"` // execute statement
	ShortCode   string     `json:"short_code"` // execute statement
	ShortURL    string     `json:"short_url"` // execute statement
	CustomAlias *string    `json:"custom_alias,omitempty"` // execute statement
	Title       *string    `json:"title,omitempty"` // execute statement
	IsActive    bool       `json:"is_active"` // execute statement
	ExpiresAt   *time.Time `json:"expires_at,omitempty"` // execute statement
	CreatedAt   time.Time  `json:"created_at"` // execute statement
	UpdatedAt   time.Time  `json:"updated_at"` // execute statement
} // end block

func NewURLResponse(url ShortURL, baseURL string) URLResponse { // declare function
	return URLResponse{ // return statement
		ID:          url.ID, // execute statement
		OriginalURL: url.OriginalURL, // execute statement
		ShortCode:   url.ShortCode, // execute statement
		ShortURL:    baseURL + "/" + url.ShortCode, // execute statement
		CustomAlias: url.CustomAlias, // execute statement
		Title:       url.Title, // execute statement
		IsActive:    url.IsActive, // execute statement
		ExpiresAt:   url.ExpiresAt, // execute statement
		CreatedAt:   url.CreatedAt, // execute statement
		UpdatedAt:   url.UpdatedAt, // execute statement
	} // end block
} // end block
