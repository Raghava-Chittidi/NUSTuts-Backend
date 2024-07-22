package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

/*
	Ping route for the cron job we set up to ensure backend is constantly being pinged.
	This makes sure that the website does not slow down automatically especially if it
	is not in use due to hosting it on render.
*/
func CronJobPingRoute(r chi.Router) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Working!"))
	})
}