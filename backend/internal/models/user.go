package models

import (
	"time"
)

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleUser    Role = "user"
	RoleOfficer Role = "officer"
)

func (r Role) Permissions() int {
	switch r {
	case RoleAdmin:
		return 3
	case RoleOfficer:
		return 2
	case RoleUser:
		return 1
	default:
		return -1
	}
}

type User struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`

	Email string `gorm:"uniqueIndex;not null"`

	Role string `gorm:"default:user"`

	Active bool `gorm:"default:true"`

	Prints []Print `gorm:"foreignKey:UserID"`
}
