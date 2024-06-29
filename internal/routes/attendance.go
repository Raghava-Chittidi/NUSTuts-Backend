package routes

import (
	"NUSTuts-Backend/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func AttendanceRoutes(r chi.Router) {
	r.Get("/{tutorialId}/list/all", handlers.GetAllAttendanceForTutorial)
	r.Get("/{tutorialId}/list", handlers.GetTodayAttendanceForTutorial)
	r.Get("/{tutorialId}", handlers.GetAttendanceCodeForTutorial)
	r.Get("/{tutorialId}/generate", handlers.GenerateAttendanceCodeForTutorial)
	r.Post("/{tutorialId}/delete", handlers.DeleteAttendanceString)
	r.Get("/student/{tutorialId}/attended/{studentId}", handlers.CheckStudentAttendance)
	r.Get("/student/{tutorialId}/attendance/{studentId}", handlers.GetStudentAttendance)
	r.Post("/student/{tutorialId}/mark", handlers.VerifyAndMarkStudentAttendance)
}