package models

import "gorm.io/gorm"

// The discussion chat for each tutorial
type Discussion struct {
	gorm.Model
	TutorialID int `json:"tutorialId" gorm:"not null;unique"`
}