package dataaccess

import (
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/models"
	"NUSTuts-Backend/internal/auth"
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

	_, claims, err := auth.AuthObj.VerifyToken(w, r)
	userID := claims.Subject // The "sub" claim is typically used for the user ID
	
	if consultation.StudentID != nil && consultation.StudentID != userID {
		return errors.New("This consultation is booked by someone else")
	}

	consultation.Booked = true
	consultation.StudentID = &userID
	database.DB.Save(&consultation)
	return nil
}

func UnbookConsultationById(id int) error {
	consultation, err := GetConsultationById(id)
	if err != nil {
		return err
	}

	_, claims, err := auth.AuthObj.VerifyToken(w, r)
	userID := claims.Subject // The "sub" claim is typically used for the user ID
	
	if consultation.StudentID != userID {
		return errors.New("You are not authorized to unbook this consultation")
	}

	consultation.Booked = false
	consultation.StudentID = nil
	database.DB.Save(&consultation)
	return nil
}