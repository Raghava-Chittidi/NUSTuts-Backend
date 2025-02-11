package util

import (
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/models"
	"fmt"
)

/*
	Generate consultations of start time and end time of 10 - 11am, 11 - 12am for a tutorial
	for every single day for the entire year
*/
func GenerateConsultationsForYear(tutorialId int, year int) error {
	/* 
		For every day in the year
		Generate a consultation in the, with date format DD-MM-YYYY
		10 - 11am
		11 - 12am 
	*/
	for i := 1; i <= 365; i++ {
		// Generate date
		date := fmt.Sprintf("%d-%d-%d", i, 1, year)
		// Generate consultation. StudentID is 0 as no student has booked the consultation
		consultation1 := models.Consultation{TutorialID: tutorialId, StudentID: 0, Date: date, StartTime: "10:00", EndTime: "11:00", Booked: false}
		consultation2 := models.Consultation{TutorialID: tutorialId, StudentID: 0, Date: date, StartTime: "11:00", EndTime: "12:00", Booked: false}
		// Save consultation
		database.DB.Create(&consultation1)
		database.DB.Create(&consultation2)
	}

	return nil
}

// Generate consultations of start time and end time of 10 - 11am, 11 - 12am for a tutorial for a given date
func GenerateConsultationsForDate(tutorialId int, date string) error {
	// Generate consultation. StudentID is 0 as no student has booked the consultation
	consultation1 := models.Consultation{TutorialID: tutorialId, StudentID: 0, Date: date, StartTime: "10:00", EndTime: "11:00", Booked: false}
	consultation2 := models.Consultation{TutorialID: tutorialId, StudentID: 0, Date: date, StartTime: "11:00", EndTime: "12:00", Booked: false}
	// Save consultation
	database.DB.Create(&consultation1)
	database.DB.Create(&consultation2)

	return nil
}


