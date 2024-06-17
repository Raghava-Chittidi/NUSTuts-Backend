package routes

import (
	"NUSTuts-Backend/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func ConsultationsRoutes(r chi.Router) {
	r.Get("/{tutorialId}", handlers.GetConsultationsForTutorialForDate)
	r.Get("/{tutorialId}/{studentId}", handlers.GetConsultationsForTutorialForStudent)
	r.Put("/{tutorialId}/book/{consultationId}", handlers.BookConsultationById)
	r.Put("/{tutorialId}/cancel/{consultationId}", handlers.CancelConsultationById)
}