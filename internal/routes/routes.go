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

// Require users to be logged in
func ProtectedRoutes() func(chi.Router) {
	return func(r chi.Router) {
		r.Route("/requests", RequestRoutes)
	}
}

/*
	Requires authorisation. Users can only perform the user actions in the tutorial they are in, 
	if they are supposed to be in that tutorial
*/
func AuthorizedRoutes() func(chi.Router) {
	return func(r chi.Router) {
		r.Route("/files", FileRoutes)
		r.Route("/messages", MessagesRoutes)
		r.Route("/ws", PrivateWebsocketRoutes)
		r.Route("/consultations", ConsultationsRoutes)
		r.Route("/attendance", AttendanceRoutes)
	}
}