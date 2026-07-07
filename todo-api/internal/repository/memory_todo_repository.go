package repository // package declaration

import ( // start import block
	"errors" // import package
	"sync" // import package

	"github.com/manasdixit190164/todo-api/internal/domain" // import package
) // end import block or close block

var ErrTodoNotFound = errors.New("todo not found") // variable declaration

type MemoryTodoRepository struct { // type/struct declaration
	mu     sync.RWMutex // statement
	todos  map[int]domain.Todo // statement
	nextID int // statement
} // statement

func NewMemoryTodoRepository() *MemoryTodoRepository { // function declaration
	return &MemoryTodoRepository{ // return result or error
		todos:  make(map[int]domain.Todo), // statement
		nextID: 1, // statement
	} // statement
} // statement

func (r *MemoryTodoRepository) Create(todo domain.Todo) (domain.Todo, error) { // function declaration
	r.mu.Lock() // statement
	defer r.mu.Unlock() // statement

	todo.ID = r.nextID // assign value
	r.todos[todo.ID] = todo // assign value
	r.nextID++ // statement

	return todo, nil // return result or error
} // statement

func (r *MemoryTodoRepository) GetAll() ([]domain.Todo, error) { // function declaration
	r.mu.RLock() // statement
	defer r.mu.RUnlock() // statement

	todos := make([]domain.Todo, 0, len(r.todos)) // declare and initialize variable

	for _, todo := range r.todos { // declare and initialize variable
		todos = append(todos, todo) // assign value
	} // statement

	return todos, nil // return result or error
} // statement

func (r *MemoryTodoRepository) GetByID(id int) (domain.Todo, error) { // function declaration
	r.mu.RLock() // statement
	defer r.mu.RUnlock() // statement

	todo, exists := r.todos[id] // declare and initialize variable
	if !exists { // check condition
		return domain.Todo{}, ErrTodoNotFound // return result or error
	} // statement

	return todo, nil // return result or error
} // statement

func (r *MemoryTodoRepository) Update(id int, todo domain.Todo) (domain.Todo, error) { // function declaration
	r.mu.Lock() // statement
	defer r.mu.Unlock() // statement

	_, exists := r.todos[id] // declare and initialize variable
	if !exists { // check condition
		return domain.Todo{}, ErrTodoNotFound // return result or error
	} // statement

	todo.ID = id // assign value
	r.todos[id] = todo // assign value

	return todo, nil // return result or error
} // statement

func (r *MemoryTodoRepository) Delete(id int) error { // function declaration
	r.mu.Lock() // statement
	defer r.mu.Unlock() // statement

	_, exists := r.todos[id] // declare and initialize variable
	if !exists { // check condition
		return ErrTodoNotFound // return result or error
	} // statement

	delete(r.todos, id) // statement
	return nil // return result or error
} // statement
