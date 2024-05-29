package models

import "gorm.io/gorm"

type Request struct {
	gorm.Model
	StudentID int `json:"studentId" gorm:"not null"`
	TutorialID int `json:"tutorialId" gorm:"not null"`
	Status string `json:"status" gorm:"not null"`
}