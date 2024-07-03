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
	r.Route("/api", func(r chi.Router) {
		r.Use(middlewares.CORS, middleware.Logger, middleware.Recoverer)
		r.Group(routes.PublicRoutes())
		r.Group(func(r chi.Router) {
			r.Use(middlewares.AuthoriseUser)
			r.Group(routes.ProtectedRoutes())
			r.Group(func(r chi.Router) {
				r.Use(middlewares.ValidateTutorialID)
				r.Group(routes.AuthorizedRoutes())
			})
		})
	})
}

func TestSetup() chi.Router {
	r := chi.NewRouter()
	setupTestRoutes(r)
	return r
}

func setupTestRoutes(r chi.Router) {
	r.Route("/api", func(r chi.Router) {
		r.Group(routes.PublicRoutes())
		r.Group(routes.ProtectedRoutes())
		r.Group(routes.AuthorizedRoutes())
	})
}
