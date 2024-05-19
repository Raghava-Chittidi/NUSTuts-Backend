package main

import (
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/util"
	"log"
)

func main()  {
	err := database.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	err = util.Migrate()
	if err != nil {
		log.Fatalln("Failed to migrate models!", err)
	}
}