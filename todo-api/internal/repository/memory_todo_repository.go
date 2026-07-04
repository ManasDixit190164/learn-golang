package repository

import (
	"errors"
	"sync"

	"github.com/manasdixit190164/todo-api/internal/domain"
)

var ErrTodoNotFound = errors.New("todo not found")

type MemoryTodoRepository struct {
	mu     sync.RWMutex
	todos  map[int]domain.Todo
	nextID int
}

func NewMemoryTodoRepository() *MemoryTodoRepository {
	return &MemoryTodoRepository{
		todos:  make(map[int]domain.Todo),
		nextID: 1,
	}
}

func (r *MemoryTodoRepository) Create(todo domain.Todo) (domain.Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	todo.ID = r.nextID
	r.todos[todo.ID] = todo
	r.nextID++

	return todo, nil
}

func (r *MemoryTodoRepository) GetAll() ([]domain.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	todos := make([]domain.Todo, 0, len(r.todos))

	for _, todo := range r.todos {
		todos = append(todos, todo)
	}

	return todos, nil
}

func (r *MemoryTodoRepository) GetByID(id int) (domain.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	todo, exists := r.todos[id]
	if !exists {
		return domain.Todo{}, ErrTodoNotFound
	}

	return todo, nil
}

func (r *MemoryTodoRepository) Update(id int, todo domain.Todo) (domain.Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.todos[id]
	if !exists {
		return domain.Todo{}, ErrTodoNotFound
	}

	todo.ID = id
	r.todos[id] = todo

	return todo, nil
}

func (r *MemoryTodoRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.todos[id]
	if !exists {
		return ErrTodoNotFound
	}

	delete(r.todos, id)
	return nil
}