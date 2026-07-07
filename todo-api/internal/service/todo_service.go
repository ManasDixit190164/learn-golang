package service // package declaration

import ( // start import block
	"errors" // import package
	"strings" // import package
	"time" // import package

	"github.com/manasdixit190164/todo-api/internal/domain" // import package
	"github.com/manasdixit190164/todo-api/internal/repository" // import package
) // end import block or close block

var ErrInvalidTodoTitle = errors.New("todo title is required") // variable declaration

type TodoService struct { // type/struct declaration
	
	repo repository.TodoRepository // statement
} // statement

func NewTodoService(repo repository.TodoRepository) *TodoService { // function declaration
	return &TodoService{ // return result or error
		repo: repo, // statement
	} // statement
} // statement

func (s *TodoService) Create(req domain.CreateTodoRequest) (domain.Todo, error) { // function declaration
	if strings.TrimSpace(req.Title) == "" { // check condition
		return domain.Todo{}, ErrInvalidTodoTitle // return result or error
	} // statement

	now := time.Now() // declare and initialize variable

	todo := domain.Todo{ // declare and initialize variable
		Title:       req.Title, // statement
		Description: req.Description, // statement
		Completed:   false, // statement
		CreatedAt:   now, // statement
		UpdatedAt:   now, // statement
	} // statement

	return s.repo.Create(todo) // return result or error
} // statement

func (s *TodoService) GetAll() ([]domain.Todo, error) { // function declaration
	return s.repo.GetAll() // return result or error
} // statement

func (s *TodoService) GetByID(id int) (domain.Todo, error) { // function declaration
	return s.repo.GetByID(id) // return result or error
} // statement

func (s *TodoService) Update(id int, req domain.UpdateTodoRequest) (domain.Todo, error) { // function declaration
	existingTodo, err := s.repo.GetByID(id) // declare and initialize variable
	if err != nil { // check condition
		return domain.Todo{}, err // return result or error
	} // statement

	if req.Title != nil { // check condition
		if strings.TrimSpace(*req.Title) == "" { // check condition
			return domain.Todo{}, ErrInvalidTodoTitle // return result or error
		} // statement
		existingTodo.Title = *req.Title // assign value
	} // statement

	if req.Description != nil { // check condition
		existingTodo.Description = *req.Description // assign value
	} // statement

	if req.Completed != nil { // check condition
		existingTodo.Completed = *req.Completed // assign value
	} // statement

	existingTodo.UpdatedAt = time.Now() // assign value

	return s.repo.Update(id, existingTodo) // return result or error
} // statement

func (s *TodoService) Delete(id int) error { // function declaration
	return s.repo.Delete(id) // return result or error
} // statement
