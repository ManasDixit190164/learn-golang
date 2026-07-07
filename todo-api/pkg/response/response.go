package response // package declaration

import ( // start import block
	"encoding/json" // import package
	"net/http" // import package
) // end import block or close block

type APIResponse struct { // type/struct declaration
	Success bool        `json:"success"` // statement
	Message string      `json:"message,omitempty"` // statement
	Data    interface{} `json:"data,omitempty"` // statement
	Error   string      `json:"error,omitempty"` // statement
} // statement

func JSON(w http.ResponseWriter, statusCode int, payload APIResponse) { // function declaration
	w.Header().Set("Content-Type", "application/json") // statement
	w.WriteHeader(statusCode) // statement
	json.NewEncoder(w).Encode(payload) // statement
} // statement
