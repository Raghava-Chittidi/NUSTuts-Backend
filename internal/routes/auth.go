package routes

import (
	"NUSTuts-Backend/internal/handlers/auth"

	"github.com/go-chi/chi/v5"
)

func AuthRoutes(r chi.Router) {
	r.Post("/signupStudent", auth.SignUpAsStudent)
	r.Post("/loginStudent", auth.LoginAsStudent)
	r.Post("/loginTA", auth.LoginAsTA)
	r.Get("/logout", auth.Logout)
}
