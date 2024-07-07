package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(isTestDb ...bool) error {
	// err := godotenv.Load("../../.env")
	// if err != nil {
	// 	log.Fatalln("Failed to load env file!")
	// 	return err
	// }
  
	var dbUrl string
	if len(isTestDb) > 0 && isTestDb[0] {
		dbUrl = os.Getenv("TEST_DB_URL")
	} else {
		dbUrl = os.Getenv("DB_URL")
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dbUrl,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		log.Fatalln("Failed to connect to database!")
		return err
	}

	DB = db
	log.Println("Connected to database!")
	return nil
}
