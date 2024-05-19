package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalln("Failed to load env file!")
		return err
	}

	db_url := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(db_url), &gorm.Config{})
	if err != nil {
		log.Fatalln("Failed to connect to database!")
		return err
	}

	DB = db
	log.Println("Connected to database!")
	return nil
}