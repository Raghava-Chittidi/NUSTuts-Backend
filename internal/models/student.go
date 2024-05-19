package models

import (
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	Name string `json:"name" gorm:"not null"`
	Email string `json:"email" gorm:"unique;not null"`
	Password string `json:"password" gorm:"not null"`
	Modules [] string `json:"modules" gorm:"type:_text;not null;"`
	Tutorials []Tutorial `json:"tutorials" gorm:"many2many:registry;type:_jsonb;not null;"`
}