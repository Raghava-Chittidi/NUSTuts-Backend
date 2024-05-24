package routes

import (
	"NUSTuts-Backend/internal/handlers/auth"
	"NUSTuts-Backend/internal/handlers/auth/studentAuth"
	"NUSTuts-Backend/internal/handlers/auth/teachingAssistantAuth"

	"github.com/go-chi/chi/v5"
)

func AuthRoutes(r chi.Router) {
	r.Post("/loginStudent", studentAuth.LoginAsStudent)
	r.Post("/loginTA", teachingAssistantAuth.LoginAsTA)
	r.Post("/signupStudent", studentAuth.SignUpAsStudent)
	r.Get("/logout", auth.Logout)
}
