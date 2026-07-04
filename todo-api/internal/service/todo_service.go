package service

import (
	"errors"
	"strings"
	"time"

	"github.com/manasdixit190164/todo-api/internal/domain"
	"github.com/manasdixit190164/todo-api/internal/repository"
)

var ErrInvalidTodoTitle = errors.New("todo title is required")

type TodoService struct {
	
	repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) *TodoService {
	return &TodoService{
		repo: repo,
	}
}

func (s *TodoService) Create(req domain.CreateTodoRequest) (domain.Todo, error) {
	if strings.TrimSpace(req.Title) == "" {
		return domain.Todo{}, ErrInvalidTodoTitle
	}

	now := time.Now()

	todo := domain.Todo{
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return s.repo.Create(todo)
}

func (s *TodoService) GetAll() ([]domain.Todo, error) {
	return s.repo.GetAll()
}

func (s *TodoService) GetByID(id int) (domain.Todo, error) {
	return s.repo.GetByID(id)
}

func (s *TodoService) Update(id int, req domain.UpdateTodoRequest) (domain.Todo, error) {
	existingTodo, err := s.repo.GetByID(id)
	if err != nil {
		return domain.Todo{}, err
	}

	if req.Title != nil {
		if strings.TrimSpace(*req.Title) == "" {
			return domain.Todo{}, ErrInvalidTodoTitle
		}
		existingTodo.Title = *req.Title
	}

	if req.Description != nil {
		existingTodo.Description = *req.Description
	}

	if req.Completed != nil {
		existingTodo.Completed = *req.Completed
	}

	existingTodo.UpdatedAt = time.Now()

	return s.repo.Update(id, existingTodo)
}

func (s *TodoService) Delete(id int) error {
	return s.repo.Delete(id)
}