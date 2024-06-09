package routes

import (
	"github.com/go-chi/chi/v5"
)

func PublicRoutes() func(chi.Router) {
	return func(r chi.Router) {
		r.Route("/auth", AuthRoutes)
	}
}

func ProtectedRoutes() func(chi.Router) {
	return func(r chi.Router) {
		r.Route("/students", StudentRoutes)
		r.Route("/teaching-assistants", TARoutes)
		r.Route("/tutorials", TutorialRoutes)
		r.Route("/requests", RequestRoutes)
		r.Route("/files", FileRoutes)
	}
}