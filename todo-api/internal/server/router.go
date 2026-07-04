package server

import (
	"net/http"

	"github.com/manasdixit190164/todo-api/internal/handler"
)

func NewRouter(todoHandler *handler.TodoHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/todos", todoHandler.TodosRoot)
	mux.HandleFunc("/todos/", todoHandler.TodosByID)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return mux
}