package domain // package declaration for the module

import ( // start import block
	"time" // import package

	"github.com/google/uuid" // import package
) // end import block or block scope

type RefreshToken struct { // declare struct type
	ID        uuid.UUID  `json:"id"` // execute statement
	UserID    uuid.UUID  `json:"user_id"` // execute statement
	TokenHash string     `json:"-"` // execute statement
	ExpiresAt time.Time  `json:"expires_at"` // execute statement
	RevokedAt *time.Time `json:"revoked_at"` // execute statement
	CreatedAt time.Time  `json:"created_at"` // execute statement
} // end block
