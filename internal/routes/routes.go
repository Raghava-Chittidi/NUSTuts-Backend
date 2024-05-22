package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func PublicRoutes() func(chi.Router) {
	return func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		
		// Logout may need to move to protected
		r.Route("/auth", AuthRoutes)
	}
}

func ProtectedRoutes() func(chi.Router) {
	return func(r chi.Router) {
		r.Route("/students", StudentRoutes)
		r.Route("/teaching-assistants", TARoutes)
		r.Route("/tutorials", TutorialRoutes)
	}
}