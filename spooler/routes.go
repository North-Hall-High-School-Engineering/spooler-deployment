package main

import (
	"context"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/torbenconto/spooler/handlers"
	"github.com/torbenconto/spooler/middleware"
	"github.com/torbenconto/spooler/models"
	"github.com/torbenconto/spooler/services"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB) (*gin.Engine, error) {
	r := gin.Default()

	// CORS config
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length", "Set-Cookie"},
		MaxAge:           12 * time.Hour,
	}))

	// Initialize services
	storageClient, err := storage.NewClient(context.Background())
	if err != nil {
		return nil, err
	}

	userSvc := services.NewUserService(db)
	otpSvc := services.NewOTPService(db)
	bucketSvc := services.NewBucketService(db, storageClient)
	printSvc := services.NewPrintService(db)

	// Public routes
	r.POST("/register", handlers.RegisterHandler(userSvc))

	otp := r.Group("/otp")
	{
		otp.POST("/request", handlers.RequestOTPHandler(otpSvc))
		otp.POST("/verify", handlers.VerifyOTPHandler(otpSvc, userSvc))
	}

	// Authenticated user routes
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/me", handlers.MeHandler())
		auth.GET("/me/prints", handlers.GetUserPrintsHandler(printSvc))
		auth.GET("/bucket/:filename", handlers.DownloadPrintFileHandler(bucketSvc))
		auth.POST("/metadata", handlers.MetadataHandler())
		auth.POST("/prints/new", handlers.NewPrintHandler(bucketSvc, printSvc))
	}

	// Admin-only routes
	admin := auth.Group("/")
	admin.Use(middleware.RoleAuthMiddleware(models.RoleAdmin))
	{
		prints := admin.Group("/prints")
		{
			prints.GET("/all", handlers.AllPrintsHandler(printSvc))

			prints.DELETE("/:id", handlers.DeletePrintHandler(printSvc, bucketSvc))
			prints.POST("/:id/approve", handlers.ApprovePrintHandler(printSvc))
			prints.POST("/:id/deny", handlers.DenyPrintHandler(printSvc))
			prints.PUT("/:id/status", handlers.UpdatePrintStatusHandler(printSvc))
		}
	}

	return r, nil
}
