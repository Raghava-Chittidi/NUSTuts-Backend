package routes

import (
	"NUSTuts-Backend/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func AuthRoutes(r chi.Router) {
	r.Post("/login", handlers.Login)
	r.Post("/signup", handlers.SignUp)
	r.Get("/logout", handlers.Logout)
}
