package models

import "gorm.io/gorm"

type EmailWhitelist struct {
	gorm.Model
	Email string `gorm:"uniqueIndex"`
}
