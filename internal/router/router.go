package router

import (
	"net/http"

	"NUSTuts-Backend/internal/routes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func StartRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/auth", routes.LoadAuthRoutes)
	router.Route("/students", routes.LoadStudentRoutes)
	router.Route("/teaching-assistants", routes.LoadTARoutes)
	router.Route("/tutorials", routes.LoadTutorialRoutes)

	return router
}
