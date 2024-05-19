package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// // Connect to db
	// err := database.Connect()
	// if err != nil {
	// 	// log.Fatalln(err)
	// 	log.Print(err)
	// }

	// // AutoMigrate db models
	// err = util.Migrate()
	// if err != nil {
	// 	// log.Fatalln("Failed to migrate models!", err)
	// 	log.Print("Failed to migrate models!", err)
	// }

	// Start server
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/hello", basicHandler)

	server := &http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	serverErr := server.ListenAndServe()
	if serverErr != nil {
		log.Fatalln(serverErr)
	}
}

func basicHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}
