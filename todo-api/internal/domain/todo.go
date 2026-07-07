package domain // package declaration

import "time" // statement

type Todo struct { // type/struct declaration
	ID          int       `json:"id"` // statement
	Title       string    `json:"title"` // statement
	Description string    `json:"description"` // statement
	Completed   bool      `json:"completed"` // statement
	CreatedAt   time.Time `json:"createdAt"` // statement
	UpdatedAt   time.Time `json:"updatedAt"` // statement
} // statement

type CreateTodoRequest struct { // type/struct declaration
	Title       string `json:"title"` // statement
	Description string `json:"description"` // statement
} // statement

type UpdateTodoRequest struct { // type/struct declaration
	Title       *string `json:"title"` // statement
	Description *string `json:"description"` // statement
	Completed   *bool   `json:"completed"` // statement
} // statement
