package repository

import "github.com/manasdixit190164/todo-api/internal/domain"

type TodoRepository interface {
	Create(todo domain.Todo) (domain.Todo, error)
	GetAll() ([]domain.Todo, error)
	GetByID(id int) (domain.Todo, error)
	Update(id int, todo domain.Todo) (domain.Todo, error)
	Delete(id int) error
}
