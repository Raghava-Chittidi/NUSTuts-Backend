package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func AllRoutes() func(r chi.Router) {
	return func (r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	
		r.Route("/auth", AuthRoutes)
		r.Route("/students", StudentRoutes)
		r.Route("/teaching-assistants", TARoutes)
		r.Route("/tutorials", TutorialRoutes)
	}
}