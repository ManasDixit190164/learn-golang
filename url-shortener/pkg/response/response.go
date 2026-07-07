package response // package declaration for the module

import ( // start import block
	"encoding/json" // import package
	"net/http" // import package
) // end import block or block scope

type APIResponse struct { // declare struct type
	Success bool        `json:"success"` // execute statement
	Message string      `json:"message,omitempty"` // execute statement
	Data    interface{} `json:"data,omitempty"` // execute statement
	Error   string      `json:"error,omitempty"` // execute statement
} // end block

func JSON(w http.ResponseWriter, statusCode int, payload APIResponse) { // declare function
	w.Header().Set("Content-Type", "application/json") // execute statement
	w.WriteHeader(statusCode) // execute statement
	_ = json.NewEncoder(w).Encode(payload) // create JSON encoder for response
} // end block

func Error(w http.ResponseWriter, statusCode int, message string) { // declare function
	JSON(w, statusCode, APIResponse{ // execute statement
		Success: false, // execute statement
		Error:   message, // execute statement
	}) // close block
} // end block
