package domain // package declaration for the module

import ( // start import block
	"time" // import package

	"github.com/google/uuid" // import package
) // end import block or block scope

type User struct { // declare struct type
	ID              uuid.UUID `json:"id"` // execute statement
	Name            string    `json:"name"` // execute statement
	Email           string    `json:"email"` // execute statement
	PasswordHash    string    `json:"-"` // execute statement
	IsEmailVerified bool      `json:"is_email_verified"` // execute statement
	CreatedAt       time.Time `json:"created_at"` // execute statement
	UpdatedAt       time.Time `json:"updated_at"` // execute statement
} // end block

type UserResponse struct { // declare struct type
	ID        uuid.UUID `json:"id"` // execute statement
	Name      string    `json:"name"` // execute statement
	Email     string    `json:"email"` // execute statement
	CreatedAt time.Time `json:"created_at"` // execute statement
} // end block

func NewUserResponse(user User) UserResponse { // declare function
	return UserResponse{ // return statement
		ID:        user.ID, // execute statement
		Name:      user.Name, // execute statement
		Email:     user.Email, // execute statement
		CreatedAt: user.CreatedAt, // execute statement
	} // end block
} // end block
