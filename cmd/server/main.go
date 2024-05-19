package main

import (
	"NUSTuts-Backend/internal/application"
	"context"
	"log"
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
	app := application.New()
	err := app.Start(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}
