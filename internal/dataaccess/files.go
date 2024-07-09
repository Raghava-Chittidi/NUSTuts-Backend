package dataaccess

import (
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/models"
	"errors"
)

func CheckIfNameExistsForTutorialIDAndWeek(tutorialId int, filename string, week int) error {
	tutorialFile, _ := GetTutorialFileFromTutorialIDAndFilename(tutorialId, filename, week)
	if tutorialFile != nil {
		return errors.New("filename already exists")
	}

	return nil
}

func CreateTutorialFile(tutorialId int, filename string, week int, filepath string) error {
	tutorialFile := &models.TutorialFile{TutorialID: tutorialId, Name: filename, Week: week, Visible: true, Filepath: filepath}
	result := database.DB.Table("tutorial_files").Create(tutorialFile)
	return result.Error
}

func GetTutorialFileFromTutorialIDAndFilename(id int, filename string, week int) (*models.TutorialFile, error) {
	var tutorialFile models.TutorialFile
	result := database.DB.Table("tutorial_files").Where("tutorial_id = ?", id).Where("week = ?", week).
				Where("name = ?", filename).First(&tutorialFile)
	
	if result.Error != nil {
		return nil, result.Error
	}

	return &tutorialFile, nil
}

func GetAllTutorialFilesFromTutorialIDAndWeek(id int, week int) (*[]models.TutorialFile, error) {
	var tutorialFiles []models.TutorialFile
	result := database.DB.Table("tutorial_files").Where("tutorial_id = ?", id).Where("week = ?", week).Find(&tutorialFiles)

	if result.Error != nil {
		return nil, result.Error
	}

	return &tutorialFiles, nil
}

func GetTutorialFileById(id int) (*models.TutorialFile, error) {
	var tutorialFile models.TutorialFile
	result := database.DB.Table("tutorial_files").Where("id = ?", id).First(&tutorialFile)

	if result.Error != nil {
		return nil, result.Error
	}

	return &tutorialFile, nil
}

func GetTutorialFileByFilepath(filepath string) (*models.TutorialFile, error) {
	var tutorialFile models.TutorialFile
	result := database.DB.Table("tutorial_files").Where("filepath = ?", filepath).First(&tutorialFile)

	if result.Error != nil {
		return nil, result.Error
	}

	return &tutorialFile, nil
}

func DeleteTutorialFileByFilepath(filepath string) (error) {
	tutorialFile, err := GetTutorialFileByFilepath(filepath)
	if err != nil {
		return err
	}

	result := database.DB.Unscoped().Table("tutorial_files").Delete(&tutorialFile)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func PrivateFileByFilepath(filepath string) error {
	tutorialFile, err := GetTutorialFileByFilepath(filepath)
	if err != nil {
		return err
	}

	tutorialFile.Visible = false
	database.DB.Save(&tutorialFile)
	return nil
}

func UnprivateFileByFilepath(filepath string) error {
	tutorialFile, err := GetTutorialFileByFilepath(filepath)
	if err != nil {
		return err
	}

	tutorialFile.Visible = true
	database.DB.Save(&tutorialFile)
	return nil
}