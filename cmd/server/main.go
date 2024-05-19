package main

import (
	"NUSTuts-Backend/internal/application"
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/util"
	"context"
	"log"
)

func main() {
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
	app := application.New()
	serverErr := app.Start(context.TODO())
	if serverErr != nil {
		log.Fatal(err)
	}
}
