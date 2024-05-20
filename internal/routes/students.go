package routes

import (
	"NUSTuts-Backend/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func StudentRoutes(r chi.Router) {
	r.Post("/", handlers.CreateStudent)
	r.Get("/", handlers.GetStudents)
	r.Get("/{id}", handlers.GetStudentByID)
	r.Put("/{id}", handlers.UpdateStudentByID)
	r.Delete("/{id}", handlers.DeleteByID)
}
