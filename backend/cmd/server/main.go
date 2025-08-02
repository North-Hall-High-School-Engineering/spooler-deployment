package main

import (
	"fmt"
	"log"

	"github.com/torbenconto/spooler/config"
	"github.com/torbenconto/spooler/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	uri := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=require",
		config.Cfg.SupabaseUser, config.Cfg.SupabasePassword, config.Cfg.SupabaseHost, config.Cfg.SupabasePort, config.Cfg.SupabaseDB,
	)

	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}

	db.AutoMigrate(models.OTP{}, models.User{}, models.Print{}, models.EmailWhitelist{})

	var admin models.User
	result := db.Where("email = ?", config.Cfg.AdminEmail).First(&admin)
	if result.Error != nil {
		admin = models.User{
			Email:     config.Cfg.AdminEmail,
			FirstName: config.Cfg.AdminFirstName,
			LastName:  config.Cfg.AdminLastName,
			Role:      string(models.RoleAdmin),
			Active:    true,
		}
		if err := db.Create(&admin).Error; err != nil {
			log.Fatalf("failed to create admin user: %v", err)
		}
	}

	// Add admin to whitelist if enabled
	if config.Cfg.EmailWhitelistEnabled {
		var count int64
		db.Model(&models.EmailWhitelist{}).Where("email = ?", config.Cfg.AdminEmail).Count(&count)
		if count == 0 {
			if err := db.Create(&models.EmailWhitelist{Email: config.Cfg.AdminEmail}).Error; err != nil {
				log.Fatalf("failed to add admin to whitelist: %v", err)
			}
		}
	}

	r, err := SetupRoutes(db)
	if err != nil {
		log.Fatalf("error setting up routes: %v", err)
	}

	if err := r.Run(fmt.Sprintf(":%s", config.Cfg.Port)); err != nil {
		log.Fatalf("failed to start server %v", err)
	}
}
