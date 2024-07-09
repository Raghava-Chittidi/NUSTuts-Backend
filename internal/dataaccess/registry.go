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

func GetAllStudentIdsOfStudentsInTutorial(tutorialId int) (*[]int, error) {
	var studentIds []int
	result := database.DB.Table("tutorials").Joins("JOIN registries ON tutorials.id = registries.tutorial_id").
				Where("tutorials.id = ?", tutorialId).Select("student_id").Find(&studentIds)
	if result.Error != nil {
		return nil, result.Error
	}

	return &studentIds, nil
}

func JoinTutorial(studentId int, tutorialId int) error {
	register := &models.Registry{StudentID: studentId, TutorialID: tutorialId}
	result := database.DB.Table("registries").Create(register)
	return result.Error
}

// func GetRegistryByStudentIDAndTutorialID(studentId int, tutorialId int) (*models.Registry, error) {
// 	var registry models.Registry
// 	result := database.DB.Table("registries").Where("student_id = ?", studentId).
// 			Where("tutorial_id = ?", tutorialId).Find(&registry)
// 	if result.Error != nil {
// 		return nil, result.Error
// 	}

// 	return &registry, nil
// }

// func DeleteRegistryByStudentIDAndTutorialID(studentId int, tutorialId int) error {
// 	registry, err := GetRegistryByStudentIDAndTutorialID(studentId, tutorialId)
// 	if err != nil {
// 		return err
// 	}

// 	result := database.DB.Unscoped().Table("registries").Delete(registry)
// 	if result.Error != nil {
// 		return result.Error
// 	}

// 	return nil
// }