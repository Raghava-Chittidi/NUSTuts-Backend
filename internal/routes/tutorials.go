package routes

import (
	"NUSTuts-Backend/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func LoadTutorialRoutes(router chi.Router) {
	tutorialsHandler := &handlers.Tutorial{}

	router.Post("/", tutorialsHandler.CreateTutorial)
	router.Get("/", tutorialsHandler.GetTutorials)
	router.Get("/{id}", tutorialsHandler.GetTutorialsByID)
	router.Put("/{id}", tutorialsHandler.UpdateTutorialByID)
}
