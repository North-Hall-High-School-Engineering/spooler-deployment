package services

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/torbenconto/spooler/models"
	"gorm.io/gorm"
)

var (
	ErrCodeInvalid = errors.New("invalid code")
	ErrCodeExpired = errors.New("code expired")
	ErrCodeUsed    = errors.New("code already used")
)

type OTPService struct {
	db *gorm.DB
}

func NewOTPService(db *gorm.DB) *OTPService {
	return &OTPService{db: db}
}

func (s *OTPService) GenerateCode(email string) (string, error) {
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	expiry := time.Now().Add(10 * time.Minute)

	otp := models.OTP{
		Email:     email,
		Code:      code,
		ExpiresAt: expiry,
		Used:      false,
	}

	if err := s.db.Create(&otp).Error; err != nil {
		return "", err
	}

	return code, nil
}

func (s *OTPService) ValidateCode(email string, code string) error {
	var otp models.OTP
	err := s.db.Where("email = ? AND code = ?", email, code).First(&otp).Error
	if err != nil {
		return ErrCodeInvalid
	}

	if otp.Used {
		return ErrCodeUsed
	}
	if otp.IsExpired() {
		return ErrCodeExpired
	}

	// Mark as used
	otp.Used = true
	if err := s.db.Save(&otp).Error; err != nil {
		return err
	}

	return nil
}
