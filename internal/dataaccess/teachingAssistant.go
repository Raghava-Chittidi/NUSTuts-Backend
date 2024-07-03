package dataaccess

import (
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/models"
)

func CreateTeachingAssistant(name string, email string, password string) (*models.TeachingAssistant, error) {
	teachingAssistant := &models.TeachingAssistant{Name: name, Email: email, Password: password}
	result := database.DB.Table("teaching_assistants").Create(teachingAssistant)
	return teachingAssistant, result.Error
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

func DeleteTeachingAssistantById(id int) error {
	teachingAssistant, err := GetTeachingAssistantById(id)
	if err != nil {
		return err
	}

	result := database.DB.Table("teaching_assistants").Delete(&teachingAssistant)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func DeleteTeachingAssistantByEmail(email string) error {
	teachingAssistant, err := GetTeachingAssistantByEmail(email)
	if err != nil {
		return err
	}

	result := database.DB.Unscoped().Table("teaching_assistants").Delete(&teachingAssistant)
	if result.Error != nil {
		return result.Error
	}

	return nil
}