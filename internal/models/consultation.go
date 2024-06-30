package models

import "gorm.io/gorm"

type Consultation struct {
	gorm.Model
	TutorialID int `json:"tutorialId" gorm:"not null uniqueIndex:idx_tutorial_id_date_start_time_end_time"`
	StudentID int `json:"studentId" gorm:"not null"`
	Date string `json:"date" gorm:"not null uniqueIndex:idx_tutorial_id_date_start_time_end_time"`
	StartTime string `json:"startTime" gorm:"not null uniqueIndex:idx_tutorial_id_date_start_time_end_time"`
	EndTime string `json:"endTime" gorm:"not null uniqueIndex:idx_tutorial_id_date_start_time_end_time"`
	Booked bool `json:"booked" gorm:"not null"`
}