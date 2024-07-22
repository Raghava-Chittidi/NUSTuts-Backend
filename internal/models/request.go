package models

import "gorm.io/gorm"

// Contains which student has made a request to which tutorial and the status of the request
type Request struct {
	gorm.Model
	StudentID int `json:"studentId" gorm:"not null"`
	TutorialID int `json:"tutorialId" gorm:"not null"`
	Status string `json:"status" gorm:"not null"` // Either accepted or rejected
}