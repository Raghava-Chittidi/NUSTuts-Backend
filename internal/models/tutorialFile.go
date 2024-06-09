package models

import "gorm.io/gorm"

type TutorialFile struct {
	gorm.Model
	TutorialID int `json:"tutorialId" gorm:"not null"`
	Filepath string `json:"filepath" gorm:"not null;unique"`
	Name string `json:"name" gorm:"not null;"`
	Visible bool `json:"visible" gorm:"not null"`
	Week int `json:"week" gorm:"not null"`
}