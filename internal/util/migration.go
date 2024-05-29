package util

import (
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/models"
)

func Migrate() error {
	err := database.DB.AutoMigrate(&models.Student{}, &models.Tutorial{}, &models.TeachingAssistant{}, &models.Request{})
	return err
}