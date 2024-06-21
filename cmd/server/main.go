package main

import (
	"NUSTuts-Backend/internal/auth"
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/router"
	"NUSTuts-Backend/internal/util"
	"NUSTuts-Backend/internal/websockets"
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

	// Initialise auth obj
	err = auth.InitialiseAuthObj()
	if err != nil {
		log.Fatalln("Failed to initialise auth obj!", err)
	}

	// Initialise and run chatrooms hub
	websockets.InitialiseHub()
	go websockets.RunHub()

	// Start server
	log.Fatalln(http.ListenAndServe(":8000", r))
}
