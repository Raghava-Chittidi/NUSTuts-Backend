package models

import "gorm.io/gorm"

type Discussion struct {
	gorm.Model
	TutorialID int `json:"tutorialId" gorm:"not null;unique"`
}