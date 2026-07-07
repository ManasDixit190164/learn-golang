package main // package declaration

import ( // start import block
	"log" // import package
	"net/http" // import package

	"github.com/manasdixit190164/todo-api/internal/config" // import package
	"github.com/manasdixit190164/todo-api/internal/handler" // import package
	"github.com/manasdixit190164/todo-api/internal/repository" // import package
	"github.com/manasdixit190164/todo-api/internal/server" // import package
	"github.com/manasdixit190164/todo-api/internal/service" // import package
) // end import block or close block

func main() { // function declaration
	cfg := config.Load() // declare and initialize variable

	todoRepo := repository.NewMemoryTodoRepository() // declare and initialize variable
	todoService := service.NewTodoService(todoRepo) // declare and initialize variable
	todoHandler := handler.NewTodoHandler(todoService) // declare and initialize variable

	router := server.NewRouter(todoHandler) // declare and initialize variable

	log.Println("server running on port", cfg.Port) // statement

	err := http.ListenAndServe(":"+cfg.Port, router) // declare and initialize variable
	if err != nil { // check condition
		log.Fatal("failed to start server:", err) // statement
	} // statement
} // statement
