package dataaccess

import (
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/models"
	"errors"
)

// DeleteConsultationById deletes a consultation by its id
func DeleteConsultationById(id int) error {
	result := database.DB.Unscoped().Table("consultations").Where("id = ?", id).Delete(&models.Consultation{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Gets all consultations for a tutorial for a specific date
func GetAllConsultationsForTutorialForDate(tutorialId int, date string) (*[]models.Consultation, error) {
	var consultations []models.Consultation
	result := database.DB.Table("consultations").
		Where("tutorial_id = ?", tutorialId).Where("date = ?", date).
		Order("date ASC").Find(&consultations)

	if result.Error != nil {
		return nil, result.Error
	}

	return &consultations, nil
}

// Gets all booked consultations booked by any student for a tutorial for a specific date for a TA
func GetBookedConsultationsForTutorialForTA(tutorialId int, date string, time string) (*[]models.Consultation, error) {
	var consultations []models.Consultation
	result := database.DB.Table("consultations").Where("tutorial_id = ?", tutorialId).
		Where("(date = ? AND end_time >= ?) OR (date > ?)", date, time, date).
		Where("booked = true").
		Order("date ASC").Order("start_time ASC").Find(&consultations)

	if result.Error != nil {
		return nil, result.Error
	}

	return &consultations, nil
}

// Gets all booked consultations booked by a student for a tutorial for a specific date
func GetBookedConsultationsForTutorialForStudent(tutorialId int, studentId int, date string, time string) (*[]models.Consultation, error) {
	var consultations []models.Consultation
	result := database.DB.Table("consultations").
		Where("tutorial_id = ?", tutorialId).Where("student_id = ?", studentId).
		Where("booked = true").
		Where("(date = ? AND end_time >= ?) OR (date > ?)", date, time, date).
		Order("date ASC").Order("start_time ASC").Find(&consultations)

	if result.Error != nil {
		return nil, result.Error
	}

	return &consultations, nil
}

// Gets consultation by its id
// Tutorial id is not needed as each consultation has a unique id
func GetConsultationById(id int) (*models.Consultation, error) {
	var consultation models.Consultation
	result := database.DB.Table("consultations").Where("id = ?", id).First(&consultation)
	if result.Error != nil {
		return nil, result.Error
	}

	return &consultation, nil
}

// Books a consultation by its id for a student with userID
func BookConsultationById(id int, userID int) (*models.Consultation, error) {
	consultation, err := GetConsultationById(id)
	if err != nil {
		return nil, err
	}

	if consultation.Booked && consultation.StudentID != userID {
		return nil, errors.New("this consultation is booked by someone else")
	}

	consultation.Booked = true
	consultation.StudentID = userID
	database.DB.Save(&consultation)
	return consultation, nil
}

// Unbooks a consultation by its id for a student with userID
func UnbookConsultationById(id int, userID int) (*models.Consultation, error) {
	consultation, err := GetConsultationById(id)
	if err != nil {
		return nil, err
	}

	if consultation.StudentID != userID {
		return nil, errors.New("you are not authorized to unbook this consultation")
	}

	consultation.Booked = false
	consultation.StudentID = 0
	database.DB.Save(&consultation)
	return consultation, nil
}
