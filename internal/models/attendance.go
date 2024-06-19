package models

import (
	"time"

	"gorm.io/gorm"
)

type Attendance struct {
	gorm.Model
	StudentID int `json:"studentId" gorm:"not null"`
	TutorialID int `json:"tutorialId" gorm:"not null"`
	Date string `json:"date" gorm:"not null"`
	Present bool `json:"present" gorm:"default:false;not null"`
}

type AttendanceString struct {
	gorm.Model
	Code string `json:"code" gorm:"not null"`
	TutorialID int `json:"tutorialId" gorm:"not null"`
	ExpiresAt time.Time `json:"expiresAt" gorm:"not null"`
}