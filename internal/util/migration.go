package util

import (
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/models"
)

// Automigrates all the models here
func Migrate() error {
	err := database.DB.AutoMigrate(&models.Student{}, &models.Tutorial{}, 
		&models.TeachingAssistant{}, &models.Request{}, &models.Registry{},
		&models.TutorialFile{}, &models.Discussion{}, &models.Message{}, 
		&models.Consultation{}, &models.Attendance{}, &models.AttendanceString{})
	return err
}