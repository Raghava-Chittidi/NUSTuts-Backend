package models

import "gorm.io/gorm"

type Consultation struct {
	gorm.Model
	TutorialID int `json:"tutorialId" gorm:"not null"`
	StudentID int `json:"studentId" gorm:"not null"`
	Date string `json:"date" gorm:"not null"`
	StartTime string `json:"startTime" gorm:"not null"`
	EndTime string `json:"endTime" gorm:"not null"`
	Booked bool `json:"booked" gorm:"not null"`
}