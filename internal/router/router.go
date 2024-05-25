package router

import (
	"NUSTuts-Backend/internal/middlewares"
	"NUSTuts-Backend/internal/routes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Setup() chi.Router {
	r := chi.NewRouter()
	setupRoutes(r)
	return r
}

func setupRoutes(r chi.Router) {
	r.Use(middlewares.CORS)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Group(routes.PublicRoutes())
		// r.Use(middlewares.AuthoriseUser)
		r.Group(routes.ProtectedRoutes())
	})
}
