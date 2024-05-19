package routes

import (
	"NUSTuts-Backend/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func LoadStudentRoutes(router chi.Router) {
	studentHandler := &handlers.Student{}

	router.Post("/", studentHandler.CreateStudent)
	router.Get("/", studentHandler.GetStudents)
	router.Get("/{id}", studentHandler.GetStudentByID)
	router.Put("/{id}", studentHandler.UpdateStudentByID)
	router.Delete("/{id}", studentHandler.DeleteByID)
}
