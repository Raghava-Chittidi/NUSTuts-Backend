package routes

import (
	"NUSTuts-Backend/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func PrivateWebsocketRoutes(r chi.Router) {
	r.Post("/{tutorialId}/create", handlers.CreateRoom)
}

func PublicWebsocketRoutes(r chi.Router) {
	r.Get("/{tutorialId}/join", handlers.JoinRoom)
}