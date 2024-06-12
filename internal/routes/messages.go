package routes

import (
	"NUSTuts-Backend/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func MessagesRoutes(r chi.Router) {
	r.Post("/", handlers.CreateMessageForTutorial)
	r.Get("/{tutorialId}", handlers.GetAllMessagesForTutorial)
}