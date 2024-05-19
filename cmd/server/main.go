package main

import (
	"log"
	"net/http"
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
	server := &http.Server{
		Addr:    ":3000",
		Handler: http.HandlerFunc(basicHandler),
	}

	serverErr := server.ListenAndServe()
	if serverErr != nil {
		log.Fatalln(serverErr)
	}
}

func basicHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}
