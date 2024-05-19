package routes

import (
	"NUSTuts-Backend/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func LoadTARoutes(router chi.Router) {
	teachingAssistantsHandler := &handlers.TeachingAssistant{}

	router.Post("/", teachingAssistantsHandler.CreateTA)
	router.Get("/", teachingAssistantsHandler.GetTAs)
	router.Get("/{id}", teachingAssistantsHandler.GetTAByID)
	router.Put("/{id}", teachingAssistantsHandler.UpdateTAByID)
	router.Delete("/{id}", teachingAssistantsHandler.DeleteTAByID)
}
