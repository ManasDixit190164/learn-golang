package handler // package declaration for the module

import ( // start import block
	"encoding/json" // import package
	"errors" // import package
	"net/http" // import package

	"github.com/manasdixit/url-shortener/internal/domain" // import package
	"github.com/manasdixit/url-shortener/internal/service" // import package
	"github.com/manasdixit/url-shortener/pkg/response" // import package
) // end import block or block scope

type AuthHandler struct { // declare struct type
	authService *service.AuthService // execute statement
} // end block

func NewAuthHandler(authService *service.AuthService) *AuthHandler { // declare function
	return &AuthHandler{authService: authService} // return statement
} // end block

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) { // declare method
	var req domain.SignupRequest // execute statement
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { // create JSON decoder for request body
		response.Error(w, http.StatusBadRequest, "invalid request body") // send an error response
		return // return statement
	} // end block

	res, err := h.authService.Signup(r.Context(), req) // declare and initialize variable
	if err != nil { // if condition
		writeAuthError(w, err) // execute statement
		return // return statement
	} // end block

	response.JSON(w, http.StatusCreated, response.APIResponse{ // send a JSON response
		Success: true, // execute statement
		Message: "signup successful", // execute statement
		Data:    res, // execute statement
	}) // close block
} // end block

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) { // declare method
	var req domain.LoginRequest // execute statement
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { // create JSON decoder for request body
		response.Error(w, http.StatusBadRequest, "invalid request body") // send an error response
		return // return statement
	} // end block

	res, err := h.authService.Login(r.Context(), req) // declare and initialize variable
	if err != nil { // if condition
		writeAuthError(w, err) // execute statement
		return // return statement
	} // end block

	response.JSON(w, http.StatusOK, response.APIResponse{ // send a JSON response
		Success: true, // execute statement
		Message: "login successful", // execute statement
		Data:    res, // execute statement
	}) // close block
} // end block

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) { // declare method
	var req domain.RefreshTokenRequest // execute statement
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { // create JSON decoder for request body
		response.Error(w, http.StatusBadRequest, "invalid request body") // send an error response
		return // return statement
	} // end block

	res, err := h.authService.Refresh(r.Context(), req.RefreshToken) // declare and initialize variable
	if err != nil { // if condition
		writeAuthError(w, err) // execute statement
		return // return statement
	} // end block

	response.JSON(w, http.StatusOK, response.APIResponse{ // send a JSON response
		Success: true, // execute statement
		Message: "token refreshed successfully", // execute statement
		Data:    res, // execute statement
	}) // close block
} // end block

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) { // declare method
	var req domain.LogoutRequest // execute statement
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { // create JSON decoder for request body
		response.Error(w, http.StatusBadRequest, "invalid request body") // send an error response
		return // return statement
	} // end block

	if err := h.authService.Logout(r.Context(), req.RefreshToken); err != nil { // if condition
		writeAuthError(w, err) // execute statement
		return // return statement
	} // end block

	response.JSON(w, http.StatusOK, response.APIResponse{ // send a JSON response
		Success: true, // execute statement
		Message: "logout successful", // execute statement
	}) // close block
} // end block

func writeAuthError(w http.ResponseWriter, err error) { // declare function
	switch { // switch statement
	case errors.Is(err, service.ErrInvalidInput): // check for a specific error
		response.Error(w, http.StatusBadRequest, "name, valid email, and password with minimum 8 characters are required") // send an error response
	case errors.Is(err, service.ErrEmailAlreadyExists): // check for a specific error
		response.Error(w, http.StatusConflict, err.Error()) // send an error response
	case errors.Is(err, service.ErrInvalidCredentials): // check for a specific error
		response.Error(w, http.StatusUnauthorized, err.Error()) // send an error response
	case errors.Is(err, service.ErrInvalidToken): // check for a specific error
		response.Error(w, http.StatusUnauthorized, err.Error()) // send an error response
	default: // execute statement
		response.Error(w, http.StatusInternalServerError, "internal server error") // send an error response
	} // end block
} // end block
