package routes

import (
	"NUSTuts-Backend/internal/handlers/auth"

	"github.com/go-chi/chi/v5"
)

func AuthRoutes(r chi.Router) {
	r.Post("/student/signup", auth.SignUpAsStudent)
	r.Post("/student/login", auth.LoginAsStudent)
	r.Post("/ta/login", auth.LoginAsTA)
	r.Get("/refresh", auth.RefreshAuthStatus)
	r.Get("/logout", auth.Logout)
}
