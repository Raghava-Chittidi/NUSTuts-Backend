package routes

import (
	"NUSTuts-Backend/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func TutorialRoutes(r chi.Router) {
	r.Post("/", handlers.CreateTutorial)
	r.Get("/", handlers.GetTutorials)
	r.Get("/{id}", handlers.GetTutorialsByID)
	r.Put("/{id}", handlers.UpdateTutorialByID)
}
