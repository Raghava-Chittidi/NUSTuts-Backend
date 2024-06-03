package routes

import (
	"NUSTuts-Backend/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func RequestRoutes(r chi.Router) {
	r.Post("/", handlers.RequestToJoinTutorial)
	r.Get("/{tutorialId}", handlers.AllPendingRequestsForTutorial)
	r.Patch("/{requestId}/accept", handlers.AcceptRequest)
	r.Patch("/{requestId}/reject", handlers.RejectRequest)
	r.Get("/{studentId}/{moduleCode}", handlers.GetUnrequestedClassNo)
}