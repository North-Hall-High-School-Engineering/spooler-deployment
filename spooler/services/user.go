package services

import (
	"errors"
	"fmt"

	"github.com/torbenconto/spooler/models"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidPIN   = errors.New("invalid pin")
	ErrEmailExists  = errors.New("email already registered")
	ErrUserInactive = errors.New("user inactive")
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// CreateUser is a raw user creation function, no validation is performed within the function itself so proper input is expected
func (s *UserService) CreateUser(user *models.User) error {
	var existingUser models.User
	if err := s.db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return ErrEmailExists
	}

	if err := s.db.Create(user).Error; err != nil {
		return err
	}

	return nil
}

// GetUserByEmail simply retrieves a user from the db by email, please keep in mind that the result of this is not to be directly served to the user without proper authentication if at all.
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}

		return nil, err
	}

	return &user, nil
}

// func (s *UserService) ValidatePIN(user *models.User, pin string) {

// }
