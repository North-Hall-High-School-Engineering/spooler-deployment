package services

import (
	"github.com/torbenconto/spooler/internal/models"
	"gorm.io/gorm"
)

type WhitelistService struct {
	db *gorm.DB
}

func NewWhitelistService(db *gorm.DB) *WhitelistService {
	return &WhitelistService{db}
}

func (s *WhitelistService) IsWhitelisted(email string) (bool, error) {
	var entry models.EmailWhitelist
	err := s.db.Where("email = ?", email).First(&entry).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	return err == nil, err
}

func (s *WhitelistService) Add(emails ...string) error {
	var entries []models.EmailWhitelist
	for _, email := range emails {
		entries = append(entries, models.EmailWhitelist{Email: email})
	}
	return s.db.Create(&entries).Error
}

func (s *WhitelistService) Remove(email string) error {
	return s.db.Where("email = ?", email).Delete(&models.EmailWhitelist{}).Error
}

func (s *WhitelistService) List() ([]models.EmailWhitelist, error) {
	var list []models.EmailWhitelist
	err := s.db.Find(&list).Error
	return list, err
}
