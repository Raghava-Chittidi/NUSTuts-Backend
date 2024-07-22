package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	Name      string         `json:"name" gorm:"not null"`
	Email     string         `json:"email" gorm:"unique;not null"`
	Password  string         `json:"password" gorm:"not null"`
	Modules   pq.StringArray `json:"modules" gorm:"type:text[];not null"` // Use pq.StringArray to prevent empty string array from being null when inserted into db
}
