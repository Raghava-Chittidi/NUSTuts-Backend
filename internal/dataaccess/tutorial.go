package dataaccess

import (
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/models"
)

func CreateTutorial(tutorialCode string, module string, teachingAssistant *models.TeachingAssistant) error {
	tutorial := &models.Tutorial{TutorialCode: tutorialCode, Module: module, TeachingAssistant: *teachingAssistant}
	result := database.DB.Table("tutorials").Create(tutorial)
	return result.Error
}

func GetTutorialById(id int) (*models.Tutorial, error) {
	var tutorial models.Tutorial
	result := database.DB.Table("tutorials").Where("id = ?", id).First(&tutorial)
	if result.Error != nil {
		return nil, result.Error
	}

	return &tutorial, nil
}