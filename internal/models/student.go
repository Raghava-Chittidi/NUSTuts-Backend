package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// A student can be in many tutorials

type Student struct {
	gorm.Model
	Name      string         `json:"name" gorm:"not null"`
	Email     string         `json:"email" gorm:"unique;not null"`
	Password  string         `json:"password" gorm:"not null"`
	Modules   pq.StringArray `json:"modules" gorm:"type:text[];not null"`              // Use pq.StringArray to prevent empty string array from being null when inserted into db
	Tutorials []Tutorial     `json:"tutorials" gorm:"many2many:registry;type:_jsonb;"` // Temporarily remove not null constraint as empty tutorial array becomes null when inserting into db
}
