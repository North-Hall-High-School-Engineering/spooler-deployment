package models

import "time"

type OTP struct {
	ID        uint      `gorm:"primaryKey"`
	Email     string    `gorm:"index;not null"`
	Code      string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"index;not null"`
	Used      bool      `gorm:"default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (otp *OTP) IsExpired() bool {
	return time.Now().After(otp.ExpiresAt)
}
