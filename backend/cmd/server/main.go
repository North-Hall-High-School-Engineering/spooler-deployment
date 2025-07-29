// @title Spooler API
// @version 1.0
// @description API for 3D print job management and authentication
// @contact.name North Hall TSA
// @contact.email northhalltsa@gmail.com
// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey Authenticated
// @in cookie
// @name token

// @securityDefinitions.apikey RoleAdmin
// @in cookie
// @name token
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

	db.AutoMigrate(models.OTP{}, models.User{}, models.Print{})

	r, err := SetupRoutes(db)
	if err != nil {
		log.Fatalf("error setting up routes: %v", err)
	}

	if err := r.Run(fmt.Sprintf(":%s", config.Cfg.Port)); err != nil {
		log.Fatalf("failed to start server %v", err)
	}
}
