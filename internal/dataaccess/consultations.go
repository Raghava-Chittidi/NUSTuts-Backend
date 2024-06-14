package dataaccess

import (
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/models"
	"errors"
)

func GetAllConsultationsForDate(tutorialId int, date string) (*[]models.Consultation, error) {
	var consultations []models.Consultation
	result := database.DB.Table("consultations")
	.Where("tutorial_id = ?", tutorialId).Where("date = ?", date).Order("date ASC").Find(&consultations)

	if result.Error != nil {
		return nil, result.Error
	}

	return &consultations, nil
}

func GetAllConsultationsForTutorialForStudent(tutorialId int, studentId int) (*[]models.Consultation, error) {
	var consultations []models.Consultation
	result := database.DB.Table("consultations")
		.Where("tutorial_id = ?", tutorialId).Where("student_id = ?", studentId)
		.Order("date ASC").Order("start_time ASC").Find(&consultations)

	if result.Error != nil {
		return nil, result.Error
	}

	return &consultations, nil
}

// Tutorial id is not needed as each consultation has a unique id
func GetConsultationById(id int) (*models.Consultation, error) {
	var request models.Consultation
	result := database.DB.Table("consultations").Where("id = ?", id).First(&request)
	if result.Error != nil {
		return nil, result.Error
	}

	return &request, nil
}

func BookConsultationById(id int) error {
	consultation, err := GetConsultationById(id)
	if err != nil {
		return err
	}

	consultation.Booked = true
	database.DB.Save(&consultation)
	return nil
}

func UnbookConsultationById(id int) error {
	consultation, err := GetConsultationById(id)
	if err != nil {
		return err
	}

	consultation.Booked = false
	database.DB.Save(&consultation)
	return nil
}