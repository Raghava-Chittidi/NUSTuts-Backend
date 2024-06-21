package routes

import (
	"NUSTuts-Backend/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func ConsultationsRoutes(r chi.Router) {
	r.Get("/{tutorialId}", handlers.GetConsultationsForTutorialForDate)
	r.Get("/student/{tutorialId}/{studentId}", handlers.GetBookedConsultationsForTutorialForStudent)
	r.Get("/teachingAssistant/{tutorialId}", handlers.GetBookedConsultationsForTutorialForTA)
	r.Put("/{tutorialId}/book/{consultationId}", handlers.BookConsultationById)
	r.Put("/{tutorialId}/cancel/{consultationId}", handlers.CancelConsultationById)
}