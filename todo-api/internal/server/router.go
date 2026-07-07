package server // package declaration

import ( // start import block
	"net/http" // import package

	"github.com/manasdixit190164/todo-api/internal/handler" // import package
) // end import block or close block

func NewRouter(todoHandler *handler.TodoHandler) *http.ServeMux { // function declaration
	mux := http.NewServeMux() // declare and initialize variable

	mux.HandleFunc("/todos", todoHandler.TodosRoot) // statement
	mux.HandleFunc("/todos/", todoHandler.TodosByID) // statement

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { // statement
		w.WriteHeader(http.StatusOK) // statement
		w.Write([]byte("OK")) // statement
	}) // statement

	return mux // return result or error
} // statement
