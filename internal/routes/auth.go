package routes

import (
	"NUSTuts-Backend/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func LoadAuthRoutes(router chi.Router) {
	authHandler := &handlers.Auth{}

	router.Post("/login", authHandler.Login)
	router.Post("/signup", authHandler.SignUp)
	router.Get("/logout", authHandler.Logout)
}
