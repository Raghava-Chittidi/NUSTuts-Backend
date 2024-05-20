package main

import (
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/router"
	"NUSTuts-Backend/internal/util"
	"log"
	"net/http"
)

func main() {
	// Setup router
	r := router.Setup()

	// Connect to db
	err := database.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	// AutoMigrate db models
	err = util.Migrate()
	if err != nil {
		log.Fatalln("Failed to migrate models!", err)
	}

	// Start server
	log.Fatalln(http.ListenAndServe(":8000", r))
}
