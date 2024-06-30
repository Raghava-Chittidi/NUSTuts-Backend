package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func CronJobPingRoute(r chi.Router) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Working!"))
	})
}