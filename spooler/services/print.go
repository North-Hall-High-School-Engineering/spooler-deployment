package services

import (
	"github.com/torbenconto/spooler/models"
	"gorm.io/gorm"
)

type PrintService struct {
	db *gorm.DB
}

func NewPrintService(db *gorm.DB) *PrintService {
	return &PrintService{
		db: db,
	}
}

func (s *PrintService) CreatePrint(print *models.Print) error {
	if err := s.db.Create(print).Error; err != nil {
		return err
	}
	return nil
}

func (s *PrintService) GetUserPrintsByID(id uint) ([]models.Print, error) {
	var prints []models.Print
	if err := s.db.Where("user_id = ?", id).Find(&prints).Error; err != nil {
		return nil, err
	}
	return prints, nil
}

func (s *PrintService) AllPrints() ([]models.Print, error) {
	var prints []models.Print
	if err := s.db.Find(&prints).Error; err != nil {
		return nil, err
	}

	return prints, nil
}

func (s *PrintService) ApprovePrint(printID uint) error {
	return s.db.Model(&models.Print{}).
		Where("id = ?", printID).
		Update("status", models.StatusPendingPrint).Error
}

func (s *PrintService) DenyPrint(printID uint, reason string) error {
	return s.db.Model(&models.Print{}).
		Where("id = ?", printID).
		Updates(map[string]interface{}{
			"status":        models.StatusDenied,
			"denial_reason": reason,
		}).Error
}

func (s *PrintService) DeletePrint(printID uint) error {
	return s.db.Delete(&models.Print{}, printID).Error
}

func (s *PrintService) GetPrintByID(id uint) (*models.Print, error) {
	var print models.Print
	if err := s.db.First(&print, id).Error; err != nil {
		return nil, err
	}
	return &print, nil
}

func (s *PrintService) UpdatePrintStatus(printID uint, status models.PrintStatus) error {
	return s.db.Model(&models.Print{}).
		Where("id = ?", printID).
		Update("status", status).Error
}
