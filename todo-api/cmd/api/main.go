package main

import (
	"log"
	"net/http"

	"github.com/manasdixit190164/todo-api/internal/config"
	"github.com/manasdixit190164/todo-api/internal/handler"
	"github.com/manasdixit190164/todo-api/internal/repository"
	"github.com/manasdixit190164/todo-api/internal/server"
	"github.com/manasdixit190164/todo-api/internal/service"
)

func main() {
	cfg := config.Load()

	todoRepo := repository.NewMemoryTodoRepository()
	todoService := service.NewTodoService(todoRepo)
	todoHandler := handler.NewTodoHandler(todoService)

	router := server.NewRouter(todoHandler)

	log.Println("server running on port", cfg.Port)

	err := http.ListenAndServe(":"+cfg.Port, router)
	if err != nil {
		log.Fatal("failed to start server:", err)
	}
}