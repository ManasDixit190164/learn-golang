package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/manasdixit190164/todo-api/internal/domain"
	"github.com/manasdixit190164/todo-api/internal/repository"
	"github.com/manasdixit190164/todo-api/internal/service"
	"github.com/manasdixit190164/todo-api/pkg/response"
)

type TodoHandler struct {
	service *service.TodoService
}

func NewTodoHandler(service *service.TodoService) *TodoHandler {
	return &TodoHandler{
		service: service,
	}
}

func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handlerMethodNotAllowed(w)
		return
	}

	var req domain.CreateTodoRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, response.APIResponse{
			Success: false,
			Error:   "invalid request body",
		})
		return
	}

	todo, err := h.service.Create(req)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	response.JSON(w, http.StatusCreated, response.APIResponse{
		Success: true,
		Message: "todo created successfully",
		Data:    todo,
	})
}

func (h *TodoHandler) GetTodos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handlerMethodNotAllowed(w)
		return
	}

	todos, err := h.service.GetAll()
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Error:   "failed to fetch todos",
		})
		return
	}

	response.JSON(w, http.StatusOK, response.APIResponse{
		Success: true,
		Data:    todos,
	})
}

func (h *TodoHandler) TodosRoot(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetTodos(w, r)
	case http.MethodPost:
		h.CreateTodo(w, r)
	default:
		handlerMethodNotAllowed(w)
	}
}

func (h *TodoHandler) TodosByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetTodoByID(w, r)
	case http.MethodPatch:
		h.UpdateTodo(w, r)
	case http.MethodDelete:
		h.DeleteTodo(w, r)
	default:
		handlerMethodNotAllowed(w)
	}
}

func (h *TodoHandler) GetTodoByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromPath(r)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.APIResponse{
			Success: false,
			Error:   "invalid todo id",
		})
		return
	}

	todo, err := h.service.GetByID(id)
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, repository.ErrTodoNotFound) {
			status = http.StatusNotFound
		}

		response.JSON(w, status, response.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	response.JSON(w, http.StatusOK, response.APIResponse{
		Success: true,
		Data:    todo,
	})
}

func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromPath(r)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.APIResponse{
			Success: false,
			Error:   "invalid todo id",
		})
		return
	}

	var req domain.UpdateTodoRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.JSON(w, http.StatusBadRequest, response.APIResponse{
			Success: false,
			Error:   "invalid request body",
		})
		return
	}

	todo, err := h.service.Update(id, req)
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, repository.ErrTodoNotFound) {
			status = http.StatusNotFound
		}

		response.JSON(w, status, response.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	response.JSON(w, http.StatusOK, response.APIResponse{
		Success: true,
		Message: "todo updated successfully",
		Data:    todo,
	})
}

func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromPath(r)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.APIResponse{
			Success: false,
			Error:   "invalid todo id",
		})
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, repository.ErrTodoNotFound) {
			status = http.StatusNotFound
		}

		response.JSON(w, status, response.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	response.JSON(w, http.StatusOK, response.APIResponse{
		Success: true,
		Message: "todo deleted successfully",
	})
}

func handlerMethodNotAllowed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("Method Not Allowed"))
}

func getIDFromPath(r *http.Request) (int, error) {
	path := strings.TrimPrefix(r.URL.Path, "/todos/")
	return strconv.Atoi(path)
}
