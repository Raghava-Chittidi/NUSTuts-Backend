package routes

import (
	"NUSTuts-Backend/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func TARoutes(r chi.Router) {
	r.Post("/", handlers.CreateTA)
	r.Get("/", handlers.GetTAs)
	r.Get("/{id}", handlers.GetTAByID)
	r.Put("/{id}", handlers.UpdateTAByID)
	r.Delete("/{id}", handlers.DeleteTAByID)
}
