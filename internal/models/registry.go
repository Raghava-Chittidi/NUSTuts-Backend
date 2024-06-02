package models

import "gorm.io/gorm"

// Contains which student is in which tutorial

type Registry struct {
	gorm.Model
	TutorialID int `json:"tutorialId" gorm:"not null"`
	StudentID int `json:"studentId" gorm:"not null"`
}