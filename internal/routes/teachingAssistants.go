package routes

import (
	"NUSTuts-Backend/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func TARoutes(r chi.Router) {
	r.Post("/", handlers.CreateTeachingAssistant)
	r.Get("/", handlers.GetTeachingAssistants)
	r.Get("/{id}", handlers.GetTeachingAssistantByID)
	r.Put("/{id}", handlers.UpdateTeachingAssistantByID)
	r.Delete("/{id}", handlers.DeleteTeachingAssistantByID)
}
