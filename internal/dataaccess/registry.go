package dataaccess

import (
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/models"
)

// Attaches student to a tutorial by creating a registry
func JoinTutorial(studentId int, tutorialId int) error {
	register := &models.Registry{StudentID: studentId, TutorialID: tutorialId}
	result := database.DB.Table("registries").Create(register)
	return result.Error
}

func GetRegistryByStudentIdAndTutorialId(studentId int, tutorialId int) (*models.Registry, error) {
	var registry models.Registry
	result := database.DB.Table("registries").Where("student_id = ?", studentId).
			Where("tutorial_id = ?", tutorialId).Find(&registry)
	if result.Error != nil {
		return nil, result.Error
	}

	return &registry, nil
}

// Get the list of the tutorials a student is in by their id
func GetTutorialsByStudentId(id int) (*[]models.Tutorial, error) {
	var tutorials []models.Tutorial
	result := database.DB.Table("tutorials").Joins("JOIN registries ON tutorials.id = registries.tutorial_id").
				Where("student_id = ?", id).Find(&tutorials)
	if result.Error != nil {
		return nil, result.Error
	}

	return &tutorials, nil
}

// Get the list of the students ids of all the students in a tutorial
func GetAllStudentIdsOfStudentsInTutorial(tutorialId int) (*[]int, error) {
	var studentIds []int
	result := database.DB.Table("tutorials").Joins("JOIN registries ON tutorials.id = registries.tutorial_id").
				Where("tutorials.id = ?", tutorialId).Select("student_id").Find(&studentIds)
	if result.Error != nil {
		return nil, result.Error
	}

	return &studentIds, nil
}

func CheckIfStudentInTutorialById(studentId int, tutorialId int) (bool, error) {
	var registry models.Registry
	result := database.DB.Table("registries").Where("tutorial_id = ?", tutorialId).Where("student_id = ?", studentId).First(&registry)
	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}

// Helper function for testing
func DeleteRegistryByStudentAndTutorial(student *models.Student, tutorial *models.Tutorial) error {
	student, err := GetStudentByEmail(student.Email)
	if err != nil {
		return err
	}

	tutorial, err = GetTutorialByClassAndModuleCode(tutorial.TutorialCode, tutorial.Module)
	if err != nil {
		return err
	}

	registry, err := GetRegistryByStudentIdAndTutorialId(int(student.ID), int(tutorial.ID))
	if err != nil {
		return err
	}

	result := database.DB.Unscoped().Table("registries").Delete(registry)
	if result.Error != nil {
		return result.Error
	}

	return nil
}