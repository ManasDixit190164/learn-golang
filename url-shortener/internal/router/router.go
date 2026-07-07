package router // package declaration for the module

import ( // start import block
	"net/http" // import package

	"github.com/go-chi/chi/v5" // import package
	"github.com/go-chi/chi/v5/middleware" // import package
	"github.com/go-chi/cors" // import package
	"github.com/manasdixit/url-shortener/internal/handler" // import package
	appmiddleware "github.com/manasdixit/url-shortener/internal/middleware" // execute statement
	"github.com/manasdixit/url-shortener/internal/utils" // import package
) // end import block or block scope

type Dependencies struct { // declare struct type
	AuthHandler *handler.AuthHandler // execute statement
	URLHandler  *handler.URLHandler // execute statement
	JWTManager  *utils.JWTManager // execute statement
} // end block

func New(deps Dependencies) http.Handler { // declare function
	r := chi.NewRouter() // declare and initialize variable

	r.Use(middleware.RequestID) // execute statement
	r.Use(middleware.RealIP) // execute statement
	r.Use(middleware.Logger) // execute statement
	r.Use(middleware.Recoverer) // execute statement
	r.Use(middleware.Heartbeat("/health")) // execute statement

	// CORS - whitelist development frontend origin and allow common methods/headers
	r.Use(cors.Handler(cors.Options{ // execute statement
		AllowedOrigins:   []string{"http: // execute statement
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}, // execute statement
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"}, // execute statement
		ExposedHeaders:   []string{"Authorization"}, // execute statement
		AllowCredentials: true, // execute statement
		MaxAge:           300, // execute statement
	})) // close block

	r.Route("/api/v1", func(r chi.Router) { // execute statement
		r.Route("/auth", func(r chi.Router) { // execute statement
			r.Post("/signup", deps.AuthHandler.Signup) // execute statement
			r.Post("/login", deps.AuthHandler.Login) // execute statement
			r.Post("/refresh", deps.AuthHandler.Refresh) // execute statement
			r.Post("/logout", deps.AuthHandler.Logout) // execute statement
		}) // close block

		r.Group(func(r chi.Router) { // execute statement
			r.Use(appmiddleware.Auth(deps.JWTManager)) // execute statement
			r.Post("/urls", deps.URLHandler.Create) // execute statement
			r.Get("/urls", deps.URLHandler.List) // execute statement
			r.Get("/urls/{id}", deps.URLHandler.GetByID) // execute statement
			r.Patch("/urls/{id}", deps.URLHandler.Update) // execute statement
			r.Delete("/urls/{id}", deps.URLHandler.Delete) // execute statement
			r.Get("/urls/{id}/analytics", deps.URLHandler.Analytics) // execute statement
		}) // close block
	}) // close block

	r.Get("/{shortCode}", deps.URLHandler.Redirect) // execute statement

	return r // return statement
} // end block
