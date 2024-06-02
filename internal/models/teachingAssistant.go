package models

import "gorm.io/gorm"

type TeachingAssistant struct {
	gorm.Model
	Name string `json:"name" gorm:"not null"`
	Email string `json:"email" gorm:"unique;not null"`
	Password string `json:"password" gorm:"not null"`
	TutorialID int `json:"tutorialId" gorm:"not null"`
}