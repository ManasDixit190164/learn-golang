package handler // package declaration for the module

import ( // start import block
	"context" // import package
	"encoding/json" // import package
	"errors" // import package
	"log/slog" // import package
	"net" // import package
	"net/http" // import package
	"strings" // import package
	"time" // import package

	"github.com/go-chi/chi/v5" // import package
	"github.com/google/uuid" // import package
	"github.com/manasdixit/url-shortener/internal/domain" // import package
	"github.com/manasdixit/url-shortener/internal/middleware" // import package
	"github.com/manasdixit/url-shortener/internal/service" // import package
	"github.com/manasdixit/url-shortener/pkg/response" // import package
) // end import block or block scope

type URLHandler struct { // declare struct type
	urlService *service.URLService // execute statement
	logger     *slog.Logger // execute statement
} // end block

func NewURLHandler(urlService *service.URLService, logger *slog.Logger) *URLHandler { // declare function
	return &URLHandler{urlService: urlService, logger: logger} // return statement
} // end block

func (h *URLHandler) Create(w http.ResponseWriter, r *http.Request) { // declare method
	userID, ok := middleware.UserIDFromContext(r.Context()) // get authenticated user ID
	if !ok { // if condition
		response.Error(w, http.StatusUnauthorized, "unauthorized") // send an error response
		return // return statement
	} // end block

	var req domain.CreateURLRequest // execute statement
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { // create JSON decoder for request body
		response.Error(w, http.StatusBadRequest, "invalid request body") // send an error response
		return // return statement
	} // end block

	res, err := h.urlService.Create(r.Context(), userID, req) // declare and initialize variable
	if err != nil { // if condition
		writeURLError(w, err) // execute statement
		return // return statement
	} // end block

	response.JSON(w, http.StatusCreated, response.APIResponse{ // send a JSON response
		Success: true, // execute statement
		Message: "short url created successfully", // execute statement
		Data:    res, // execute statement
	}) // close block
} // end block

func (h *URLHandler) List(w http.ResponseWriter, r *http.Request) { // declare method
	userID, ok := middleware.UserIDFromContext(r.Context()) // get authenticated user ID
	if !ok { // if condition
		response.Error(w, http.StatusUnauthorized, "unauthorized") // send an error response
		return // return statement
	} // end block

	res, err := h.urlService.List(r.Context(), userID) // declare and initialize variable
	if err != nil { // if condition
		response.Error(w, http.StatusInternalServerError, "internal server error") // send an error response
		return // return statement
	} // end block

	response.JSON(w, http.StatusOK, response.APIResponse{Success: true, Data: res}) // send a JSON response
} // end block

func (h *URLHandler) GetByID(w http.ResponseWriter, r *http.Request) { // declare method
	userID, ok := middleware.UserIDFromContext(r.Context()) // get authenticated user ID
	if !ok { // if condition
		response.Error(w, http.StatusUnauthorized, "unauthorized") // send an error response
		return // return statement
	} // end block

	id, err := parseUUIDParam(r, "id") // declare and initialize variable
	if err != nil { // if condition
		response.Error(w, http.StatusBadRequest, "invalid url id") // send an error response
		return // return statement
	} // end block

	res, err := h.urlService.GetByID(r.Context(), userID, id) // declare and initialize variable
	if err != nil { // if condition
		writeURLError(w, err) // execute statement
		return // return statement
	} // end block

	response.JSON(w, http.StatusOK, response.APIResponse{Success: true, Data: res}) // send a JSON response
} // end block

func (h *URLHandler) Update(w http.ResponseWriter, r *http.Request) { // declare method
	userID, ok := middleware.UserIDFromContext(r.Context()) // get authenticated user ID
	if !ok { // if condition
		response.Error(w, http.StatusUnauthorized, "unauthorized") // send an error response
		return // return statement
	} // end block

	id, err := parseUUIDParam(r, "id") // declare and initialize variable
	if err != nil { // if condition
		response.Error(w, http.StatusBadRequest, "invalid url id") // send an error response
		return // return statement
	} // end block

	var req domain.UpdateURLRequest // execute statement
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { // create JSON decoder for request body
		response.Error(w, http.StatusBadRequest, "invalid request body") // send an error response
		return // return statement
	} // end block

	res, err := h.urlService.Update(r.Context(), userID, id, req) // declare and initialize variable
	if err != nil { // if condition
		writeURLError(w, err) // execute statement
		return // return statement
	} // end block

	response.JSON(w, http.StatusOK, response.APIResponse{ // send a JSON response
		Success: true, // execute statement
		Message: "url updated successfully", // execute statement
		Data:    res, // execute statement
	}) // close block
} // end block

func (h *URLHandler) Delete(w http.ResponseWriter, r *http.Request) { // declare method
	userID, ok := middleware.UserIDFromContext(r.Context()) // get authenticated user ID
	if !ok { // if condition
		response.Error(w, http.StatusUnauthorized, "unauthorized") // send an error response
		return // return statement
	} // end block

	id, err := parseUUIDParam(r, "id") // declare and initialize variable
	if err != nil { // if condition
		response.Error(w, http.StatusBadRequest, "invalid url id") // send an error response
		return // return statement
	} // end block

	if err := h.urlService.Delete(r.Context(), userID, id); err != nil { // if condition
		writeURLError(w, err) // execute statement
		return // return statement
	} // end block

	response.JSON(w, http.StatusOK, response.APIResponse{ // send a JSON response
		Success: true, // execute statement
		Message: "url deactivated successfully", // execute statement
	}) // close block
} // end block

func (h *URLHandler) Analytics(w http.ResponseWriter, r *http.Request) { // declare method
	userID, ok := middleware.UserIDFromContext(r.Context()) // get authenticated user ID
	if !ok { // if condition
		response.Error(w, http.StatusUnauthorized, "unauthorized") // send an error response
		return // return statement
	} // end block

	id, err := parseUUIDParam(r, "id") // declare and initialize variable
	if err != nil { // if condition
		response.Error(w, http.StatusBadRequest, "invalid url id") // send an error response
		return // return statement
	} // end block

	res, err := h.urlService.Analytics(r.Context(), userID, id) // declare and initialize variable
	if err != nil { // if condition
		writeURLError(w, err) // execute statement
		return // return statement
	} // end block

	response.JSON(w, http.StatusOK, response.APIResponse{Success: true, Data: res}) // send a JSON response
} // end block

func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) { // declare method
	shortCode := chi.URLParam(r, "shortCode") // read URL path parameter
	if strings.TrimSpace(shortCode) == "" { // trim whitespace
		response.Error(w, http.StatusBadRequest, "invalid short code") // send an error response
		return // return statement
	} // end block

	resolvedURL, err := h.urlService.ResolveRedirect(r.Context(), shortCode) // declare and initialize variable
	if err != nil { // if condition
		writeURLError(w, err) // execute statement
		return // return statement
	} // end block

	click := domain.Click{ // declare and initialize variable
		URLID:     resolvedURL.ID, // execute statement
		IPAddress: clientIP(r), // execute statement
		UserAgent: r.UserAgent(), // execute statement
		Referrer:  r.Referer(), // execute statement
	} // end block

	go func() { // start goroutine
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // create a context with timeout
		defer cancel() // defer function call
		if err := h.urlService.TrackClick(ctx, click); err != nil { // if condition
			h.logger.Warn("failed to track click", "error", err, "url_id", resolvedURL.ID) // execute statement
		} // end block
	}() // close block

	http.Redirect(w, r, resolvedURL.OriginalURL, http.StatusFound) // redirect client to destination URL
} // end block

func parseUUIDParam(r *http.Request, key string) (uuid.UUID, error) { // declare function
	value := chi.URLParam(r, key) // read URL path parameter
	return uuid.Parse(value) // return statement
} // end block

func clientIP(r *http.Request) string { // declare function
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" { // if condition
		parts := strings.Split(ip, ",") // declare and initialize variable
		return strings.TrimSpace(parts[0]) // trim whitespace
	} // end block
	if ip := r.Header.Get("X-Real-IP"); ip != "" { // if condition
		return ip // return statement
	} // end block
	host, _, err := net.SplitHostPort(r.RemoteAddr) // declare and initialize variable
	if err != nil { // if condition
		return r.RemoteAddr // return statement
	} // end block
	return host // return statement
} // end block

func writeURLError(w http.ResponseWriter, err error) { // declare function
	switch { // switch statement
	case errors.Is(err, service.ErrInvalidURL): // check for a specific error
		response.Error(w, http.StatusBadRequest, err.Error()) // send an error response
	case errors.Is(err, service.ErrInvalidInput): // check for a specific error
		response.Error(w, http.StatusBadRequest, err.Error()) // send an error response
	case errors.Is(err, service.ErrAliasUnavailable): // check for a specific error
		response.Error(w, http.StatusConflict, err.Error()) // send an error response
	case errors.Is(err, service.ErrURLNotFound): // check for a specific error
		response.Error(w, http.StatusNotFound, err.Error()) // send an error response
	case errors.Is(err, service.ErrURLInactive): // check for a specific error
		response.Error(w, http.StatusGone, err.Error()) // send an error response
	case errors.Is(err, service.ErrURLExpired): // check for a specific error
		response.Error(w, http.StatusGone, err.Error()) // send an error response
	default: // execute statement
		response.Error(w, http.StatusInternalServerError, "internal server error") // send an error response
	} // end block
} // end block
