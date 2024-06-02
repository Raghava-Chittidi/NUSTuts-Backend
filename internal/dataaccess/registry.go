package dataaccess

import (
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/models"
)

func GetTutorialsByStudentId(id int) (*[]models.Tutorial, error) {
	var tutorials []models.Tutorial
	result := database.DB.Table("tutorials").Joins("JOIN registries ON tutorials.id = registries.tutorial_id").
				Where("student_id = ?", id).Find(&tutorials)
	if result.Error != nil {
		return nil, result.Error
	}

	return &tutorials, nil
}

func JoinTutorial(studentId int, tutorialId int) error {
	register := &models.Registry{StudentID: studentId, TutorialID: tutorialId}
	result := database.DB.Table("registries").Create(register)
	return result.Error
}