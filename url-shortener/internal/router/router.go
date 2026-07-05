package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/manasdixit/url-shortener/internal/handler"
	appmiddleware "github.com/manasdixit/url-shortener/internal/middleware"
	"github.com/manasdixit/url-shortener/internal/utils"
)

type Dependencies struct {
	AuthHandler *handler.AuthHandler
	URLHandler  *handler.URLHandler
	JWTManager  *utils.JWTManager
}

func New(deps Dependencies) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/health"))

	// CORS - whitelist development frontend origin and allow common methods/headers
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Authorization"},
		AllowCredentials: true,
		MaxAge:           300, // 5 minutes
	}))

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/signup", deps.AuthHandler.Signup)
			r.Post("/login", deps.AuthHandler.Login)
			r.Post("/refresh", deps.AuthHandler.Refresh)
			r.Post("/logout", deps.AuthHandler.Logout)
		})

		r.Group(func(r chi.Router) {
			r.Use(appmiddleware.Auth(deps.JWTManager))
			r.Post("/urls", deps.URLHandler.Create)
			r.Get("/urls", deps.URLHandler.List)
			r.Get("/urls/{id}", deps.URLHandler.GetByID)
			r.Patch("/urls/{id}", deps.URLHandler.Update)
			r.Delete("/urls/{id}", deps.URLHandler.Delete)
			r.Get("/urls/{id}/analytics", deps.URLHandler.Analytics)
		})
	})

	r.Get("/{shortCode}", deps.URLHandler.Redirect)

	return r
}
