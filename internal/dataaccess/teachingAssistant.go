package dataaccess

import (
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/models"
)

func CreateTeachingAssistant(name string, email string, password string, tutorialId uint) error {
	teachingAssistant := &models.TeachingAssistant{Name: name, Email: email, Password: password, TutorialID: tutorialId}
	result := database.DB.Table("teaching_assistants").Create(teachingAssistant)
	return result.Error
}

func GetTeachingAssistantById(id int) (*models.TeachingAssistant, error) {
	var teachingAssistant models.TeachingAssistant
	result := database.DB.Table("teaching_assistants").Where("id = ?", id).First(&teachingAssistant)
	if result.Error != nil {
		return nil, result.Error
	}

	return &teachingAssistant, nil
}

func GetTeachingAssistantByEmail(email string) (*models.TeachingAssistant, error) {
	var teachingAssistant models.TeachingAssistant
	result := database.DB.Table("teaching_assistants").Where("email = ?", email).First(&teachingAssistant)
	if result.Error != nil {
		return nil, result.Error
	}

	return &teachingAssistant, nil
}