package models

import "gorm.io/gorm"

type Tutorial struct {
	gorm.Model
	TutorialCode string `json:"tutorialCode" gorm:"not null"`
	Module string `json:"module" gorm:"not null"`
	TeachingAssistant TeachingAssistant `json:"teachingAssistant" gorm:"type:jsonb;not null"`
	Students []Student `json:"students" gorm:"many2many:registry;type:_jsonb;not null;"`
}