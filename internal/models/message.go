package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	DiscussionID int `json:"discussionId" gorm:"not null"`
	SenderID int `json:"senderId" gorm:"not null"`
	UserType string `json:"userType" gorm:"not null"`
	Content string `json:"content" gorm:"not null"`
}