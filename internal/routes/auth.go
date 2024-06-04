package routes

import (
	"NUSTuts-Backend/internal/handlers/auth"

	"github.com/go-chi/chi/v5"
)

func AuthRoutes(r chi.Router) {
	r.Post("/students/signup", auth.SignUpAsStudent)
	r.Post("/students/login", auth.LoginAsStudent)
	r.Post("/teaching-assistants/login", auth.LoginAsTeachingAssistant)
	r.Get("/logout", auth.Logout)
}
