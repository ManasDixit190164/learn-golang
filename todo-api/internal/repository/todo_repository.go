package repository // package declaration

import "github.com/manasdixit190164/todo-api/internal/domain" // statement

type TodoRepository interface { // type declaration
	Create(todo domain.Todo) (domain.Todo, error)         // statement
	GetAll() ([]domain.Todo, error)                       // statement
	GetByID(id int) (domain.Todo, error)                  // statement
	Update(id int, todo domain.Todo) (domain.Todo, error) // statement
	Delete(id int) error                                  // declare function signature
} // statement
