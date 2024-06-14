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
		r.Route("/student", StudentRoutes)
		r.Route("/teachingAssistant", TARoutes)
		r.Route("/tutorials", TutorialRoutes)
		r.Route("/requests", RequestRoutes)
	}
}

func AuthorizedRoutes() func(chi.Router) {
	return func(r chi.Router) {
		r.Route("/files", FileRoutes)
		r.Route("/messages", MessagesRoutes)
		r.Route("/consultations", ConsultationsRoutes)
	}
}