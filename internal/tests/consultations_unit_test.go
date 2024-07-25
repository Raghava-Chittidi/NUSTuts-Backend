package tests

import (
	"NUSTuts-Backend/internal/dataaccess"
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/models"
	"testing"

	"NUSTuts-Backend/internal/util"

	"github.com/stretchr/testify/assert"
)

func TestGenerateConsultationsForDate(t *testing.T) {
	// Make sure consultations table is empty
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Consultation{})

	// Current no. of consultations in the test db should be 0
	var count int64
	database.DB.Table("consultations").Count(&count)
	assert.Equal(t, 0, int(count))

	// Generate consultations for the date
	err := util.GenerateConsultationsForDate(1, "2021-01-01")
	assert.NoError(t, err)

	// Current no. of consultations in the test db should be 2
	database.DB.Table("consultations").Count(&count)
	assert.Equal(t, 2, int(count))

	// Check if the consultations were generated properly
	var consultations []models.Consultation
	database.DB.Find(&consultations)
	assert.Equal(t, 2, len(consultations))
	assert.Equal(t, 1, consultations[0].TutorialID)
	assert.Equal(t, 1, consultations[1].TutorialID)
	assert.Equal(t, 0, consultations[0].StudentID)
	assert.Equal(t, 0, consultations[1].StudentID)
	assert.Equal(t, "2021-01-01", consultations[0].Date)
	assert.Equal(t, "2021-01-01", consultations[1].Date)
	assert.Equal(t, "10:00", consultations[0].StartTime)
	assert.Equal(t, "11:00", consultations[0].EndTime)
	assert.Equal(t, "11:00", consultations[1].StartTime)
	assert.Equal(t, "12:00", consultations[1].EndTime)
	assert.False(t, consultations[0].Booked)
	assert.False(t, consultations[1].Booked)

	// Cleanup
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Consultation{})
}

func TestGetConsultationById(t *testing.T) {
	// Make sure consultations table is empty
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Consultation{})

	// Generate consultations for the date
	err := util.GenerateConsultationsForDate(1, "2021-01-01")
	assert.NoError(t, err)

	// Get consultation generated for timeslot 10:00 - 11:00
	var consultation models.Consultation
	database.DB.Where("tutorial_id = ? AND date = ? AND start_time = ?", 1, "2021-01-01", "10:00").First(&consultation)

	// Get the consultation by ID
	consultationByID, err := dataaccess.GetConsultationById(int(consultation.ID))
	assert.NoError(t, err)

	// Check if the consultations are the same
	assert.Equal(t, consultation.ID, consultationByID.ID)
	assert.Equal(t, consultation.TutorialID, consultationByID.TutorialID)
	assert.Equal(t, consultation.StudentID, consultationByID.StudentID)
	assert.Equal(t, consultation.Date, consultationByID.Date)
	assert.Equal(t, consultation.StartTime, consultationByID.StartTime)
	assert.Equal(t, consultation.EndTime, consultationByID.EndTime)
	assert.Equal(t, consultation.Booked, consultationByID.Booked)

	// Cleanup
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Consultation{})
}

func TestBookUnbookedConsultation(t *testing.T) {
	// Make sure consultations table is empty
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Consultation{})

	// Generate consultations for the date
	err := util.GenerateConsultationsForDate(1, "2021-01-01")
	assert.NoError(t, err)

	// Get consultation generated for timeslot 10:00 - 11:00
	var consultation models.Consultation
	database.DB.Where("tutorial_id = ? AND date = ? AND start_time = ?", 1, "2021-01-01", "10:00").First(&consultation)

	// Book a consultation
	_, err = dataaccess.BookConsultationById(int(consultation.ID), 1)
	assert.NoError(t, err)

	// Check if the consultation was booked properly
	database.DB.Where("tutorial_id = ? AND student_id = ? AND date = ? AND start_time = ?", 1, 1, "2021-01-01", "10:00").First(&consultation)
	assert.True(t, consultation.Booked)

	// Cleanup
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Consultation{})
}

func TestBookBookedConsultation(t *testing.T) {
	// Make sure consultations table is empty
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Consultation{})

	// Generate consultations for the date
	err := util.GenerateConsultationsForDate(1, "2021-01-01")
	assert.NoError(t, err)

	// Get consultation generated for timeslot 10:00 - 11:00
	var consultation models.Consultation
	database.DB.Where("tutorial_id = ? AND date = ? AND start_time = ?", 1, "2021-01-01", "10:00").First(&consultation)

	// Book a consultation
	_, err = dataaccess.BookConsultationById(int(consultation.ID), 1)
	assert.NoError(t, err)

	// Try booking the same consultation again with a different student ID
	_, err = dataaccess.BookConsultationById(int(consultation.ID), 2)
	assert.Error(t, err)
	// Assert error message
	assert.Equal(t, "this consultation is booked by someone else", err.Error())

	// Cleanup
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Consultation{})
}

func TestValidUnbookBookedConsultation(t *testing.T) {
	// Make sure consultations table is empty
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Consultation{})

	// Generate consultations for the date
	err := util.GenerateConsultationsForDate(1, "2021-01-01")
	assert.NoError(t, err)

	// Get consultation generated for timeslot 10:00 - 11:00
	var consultation models.Consultation
	database.DB.Where("tutorial_id = ? AND date = ? AND start_time = ?", 1, "2021-01-01", "10:00").First(&consultation)

	// Book a consultation
	_, err = dataaccess.BookConsultationById(int(consultation.ID), 1)
	assert.NoError(t, err)

	// Unbook the consultation
	_, err = dataaccess.UnbookConsultationById(int(consultation.ID), 1)
	assert.NoError(t, err)

	// Check if the consultation was unbooked properly
	database.DB.Where("tutorial_id = ? AND student_id = ? AND date = ? AND start_time = ?", 1, 0, "2021-01-01", "10:00").First(&consultation)
	assert.False(t, consultation.Booked)

	// Cleanup
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Consultation{})
}

func TestInvalidUnbookBookedConsultation(t *testing.T) {
	// Make sure consultations table is empty
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Consultation{})

	// Generate consultations for the date
	err := util.GenerateConsultationsForDate(1, "2021-01-01")
	assert.NoError(t, err)

	// Get consultation generated for timeslot 10:00 - 11:00
	var consultation models.Consultation
	database.DB.Where("tutorial_id = ? AND date = ? AND start_time = ?", 1, "2021-01-01", "10:00").First(&consultation)

	// Book a consultation
	_, err = dataaccess.BookConsultationById(int(consultation.ID), 1)
	assert.NoError(t, err)

	// Try unbooking the same consultation again with a different student ID
	_, err = dataaccess.UnbookConsultationById(int(consultation.ID), 2)
	assert.Error(t, err)
	// Assert error message
	assert.Equal(t, "you are not authorized to unbook this consultation", err.Error())

	// Cleanup
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Consultation{})
}

func TestDeleteConsultationById(t *testing.T) {
	// Make sure consultations table is empty
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Consultation{})

	// Generate consultations for the date
	err := util.GenerateConsultationsForDate(1, "2021-01-01")
	assert.NoError(t, err)

	// Get consultation generated for timeslot 10:00 - 11:00
	var consultation models.Consultation
	database.DB.Where("tutorial_id = ? AND date = ? AND start_time = ?", 1, "2021-01-01", "10:00").First(&consultation)

	// Delete the consultation
	err = dataaccess.DeleteConsultationById(int(consultation.ID))
	assert.NoError(t, err)

	// Check if the consultation was deleted properly
	result := database.DB.Where("tutorial_id = ? AND date = ? AND start_time = ?", 1, "2021-01-01", "10:00").First(&consultation)
	assert.Error(t, result.Error)

	// Cleanup
	database.DB.Unscoped().Where("1 = 1").Delete(&models.Consultation{})
}
