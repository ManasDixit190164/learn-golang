package handler // package declaration

import ( // start import block
	"encoding/json" // import package
	"errors" // import package
	"net/http" // import package
	"strconv" // import package
	"strings" // import package

	"github.com/manasdixit190164/todo-api/internal/domain" // import package
	"github.com/manasdixit190164/todo-api/internal/repository" // import package
	"github.com/manasdixit190164/todo-api/internal/service" // import package
	"github.com/manasdixit190164/todo-api/pkg/response" // import package
) // end import block or close block

type TodoHandler struct { // type/struct declaration
	service *service.TodoService // statement
} // statement

func NewTodoHandler(service *service.TodoService) *TodoHandler { // function declaration
	return &TodoHandler{ // return result or error
		service: service, // statement
	} // statement
} // statement

func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) { // function declaration
	if r.Method != http.MethodPost { // check condition
		handlerMethodNotAllowed(w) // statement
		return // return
	} // statement

	var req domain.CreateTodoRequest // variable declaration

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { // decode JSON request into struct
		response.JSON(w, http.StatusBadRequest, response.APIResponse{ // write a JSON API response
			Success: false, // statement
			Error:   "invalid request body", // statement
		}) // statement
		return // return
	} // statement

	todo, err := h.service.Create(req) // declare and initialize variable
	if err != nil { // check condition
		response.JSON(w, http.StatusBadRequest, response.APIResponse{ // write a JSON API response
			Success: false, // statement
			Error:   err.Error(), // statement
		}) // statement
		return // return
	} // statement

	response.JSON(w, http.StatusCreated, response.APIResponse{ // write a JSON API response
		Success: true, // statement
		Message: "todo created successfully", // statement
		Data:    todo, // statement
	}) // statement
} // statement

func (h *TodoHandler) GetTodos(w http.ResponseWriter, r *http.Request) { // function declaration
	if r.Method != http.MethodGet { // check condition
		handlerMethodNotAllowed(w) // statement
		return // return
	} // statement

	todos, err := h.service.GetAll() // declare and initialize variable
	if err != nil { // check condition
		response.JSON(w, http.StatusInternalServerError, response.APIResponse{ // write a JSON API response
			Success: false, // statement
			Error:   "failed to fetch todos", // statement
		}) // statement
		return // return
	} // statement

	response.JSON(w, http.StatusOK, response.APIResponse{ // write a JSON API response
		Success: true, // statement
		Data:    todos, // statement
	}) // statement
} // statement

func (h *TodoHandler) TodosRoot(w http.ResponseWriter, r *http.Request) { // function declaration
	switch r.Method { // statement
	case http.MethodGet: // statement
		h.GetTodos(w, r) // statement
	case http.MethodPost: // statement
		h.CreateTodo(w, r) // statement
	default: // statement
		handlerMethodNotAllowed(w) // statement
	} // statement
} // statement

func (h *TodoHandler) TodosByID(w http.ResponseWriter, r *http.Request) { // function declaration
	switch r.Method { // statement
	case http.MethodGet: // statement
		h.GetTodoByID(w, r) // statement
	case http.MethodPatch: // statement
		h.UpdateTodo(w, r) // statement
	case http.MethodDelete: // statement
		h.DeleteTodo(w, r) // statement
	default: // statement
		handlerMethodNotAllowed(w) // statement
	} // statement
} // statement

func (h *TodoHandler) GetTodoByID(w http.ResponseWriter, r *http.Request) { // function declaration
	id, err := getIDFromPath(r) // declare and initialize variable
	if err != nil { // check condition
		response.JSON(w, http.StatusBadRequest, response.APIResponse{ // write a JSON API response
			Success: false, // statement
			Error:   "invalid todo id", // statement
		}) // statement
		return // return
	} // statement

	todo, err := h.service.GetByID(id) // declare and initialize variable
	if err != nil { // check condition
		status := http.StatusInternalServerError // declare and initialize variable

		if errors.Is(err, repository.ErrTodoNotFound) { // check condition
			status = http.StatusNotFound // assign value
		} // statement

		response.JSON(w, status, response.APIResponse{ // write a JSON API response
			Success: false, // statement
			Error:   err.Error(), // statement
		}) // statement
		return // return
	} // statement

	response.JSON(w, http.StatusOK, response.APIResponse{ // write a JSON API response
		Success: true, // statement
		Data:    todo, // statement
	}) // statement
} // statement

func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) { // function declaration
	id, err := getIDFromPath(r) // declare and initialize variable
	if err != nil { // check condition
		response.JSON(w, http.StatusBadRequest, response.APIResponse{ // write a JSON API response
			Success: false, // statement
			Error:   "invalid todo id", // statement
		}) // statement
		return // return
	} // statement

	var req domain.UpdateTodoRequest // variable declaration

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { // decode JSON request into struct
		response.JSON(w, http.StatusBadRequest, response.APIResponse{ // write a JSON API response
			Success: false, // statement
			Error:   "invalid request body", // statement
		}) // statement
		return // return
	} // statement

	todo, err := h.service.Update(id, req) // declare and initialize variable
	if err != nil { // check condition
		status := http.StatusInternalServerError // declare and initialize variable

		if errors.Is(err, repository.ErrTodoNotFound) { // check condition
			status = http.StatusNotFound // assign value
		} // statement

		response.JSON(w, status, response.APIResponse{ // write a JSON API response
			Success: false, // statement
			Error:   err.Error(), // statement
		}) // statement
		return // return
	} // statement

	response.JSON(w, http.StatusOK, response.APIResponse{ // write a JSON API response
		Success: true, // statement
		Message: "todo updated successfully", // statement
		Data:    todo, // statement
	}) // statement
} // statement

func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) { // function declaration
	id, err := getIDFromPath(r) // declare and initialize variable
	if err != nil { // check condition
		response.JSON(w, http.StatusBadRequest, response.APIResponse{ // write a JSON API response
			Success: false, // statement
			Error:   "invalid todo id", // statement
		}) // statement
		return // return
	} // statement

	err = h.service.Delete(id) // assign value
	if err != nil { // check condition
		status := http.StatusInternalServerError // declare and initialize variable

		if errors.Is(err, repository.ErrTodoNotFound) { // check condition
			status = http.StatusNotFound // assign value
		} // statement

		response.JSON(w, status, response.APIResponse{ // write a JSON API response
			Success: false, // statement
			Error:   err.Error(), // statement
		}) // statement
		return // return
	} // statement

	response.JSON(w, http.StatusOK, response.APIResponse{ // write a JSON API response
		Success: true, // statement
		Message: "todo deleted successfully", // statement
	}) // statement
} // statement

func handlerMethodNotAllowed(w http.ResponseWriter) { // function declaration
	w.WriteHeader(http.StatusMethodNotAllowed) // statement
	w.Write([]byte("Method Not Allowed")) // statement
} // statement

func getIDFromPath(r *http.Request) (int, error) { // function declaration
	path := strings.TrimPrefix(r.URL.Path, "/todos/") // declare and initialize variable
	return strconv.Atoi(path) // return result or error
} // statement
