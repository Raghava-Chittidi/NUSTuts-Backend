package routes

import (
	"github.com/go-chi/chi/v5"
)

func PublicRoutes() func(chi.Router) {
	return func(r chi.Router) {
		r.Route("/auth", AuthRoutes)
		r.Route("/public/ws", PublicWebsocketRoutes)
		r.Route("/ping", CronJobPingRoute)
	}
}

func ProtectedRoutes() func(chi.Router) {
	return func(r chi.Router) {
		r.Route("/requests", RequestRoutes)
	}
}

func AuthorizedRoutes() func(chi.Router) {
	return func(r chi.Router) {
		r.Route("/files", FileRoutes)
		r.Route("/messages", MessagesRoutes)
		r.Route("/ws", PrivateWebsocketRoutes)
		r.Route("/consultations", ConsultationsRoutes)
		r.Route("/attendance", AttendanceRoutes)
	}
}