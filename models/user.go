package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Email    string `gorm:"not null;uniqueIndex"`
	Password string `gorm:"not null"`
}
