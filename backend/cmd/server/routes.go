package main

import (
	"context"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/torbenconto/spooler/internal/handlers"
	"github.com/torbenconto/spooler/internal/middleware"
	"github.com/torbenconto/spooler/internal/models"
	"github.com/torbenconto/spooler/internal/services"
	"gorm.io/gorm"
)

func getAllowedOrigins() []string {
	origins := os.Getenv("CORS_ALLOW_ORIGINS")
	if origins == "" {
		return []string{"http://localhost:5173"}
	}
	parts := strings.Split(origins, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func SetupRoutes(db *gorm.DB) (*gin.Engine, error) {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     getAllowedOrigins(),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length", "Set-Cookie"},
		MaxAge:           12 * time.Hour,
	}))

	storageClient, err := storage.NewClient(context.Background())
	if err != nil {
		return nil, err
	}

	userSvc := services.NewUserService(db)
	otpSvc := services.NewOTPService(db)
	bucketSvc := services.NewBucketService(db, storageClient)
	printSvc := services.NewPrintService(db)
	whitelistSvc := services.NewWhitelistService(db)

	// Public routes
	otp := r.Group("/otp")
	{
		otp.POST("/request", handlers.RequestOTPHandler(otpSvc, whitelistSvc))
		otp.POST("/verify", handlers.VerifyOTPHandler(otpSvc, userSvc))
	}
	r.POST("/register", handlers.RegisterHandler(userSvc, whitelistSvc))

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
			prints.PUT("/:id", handlers.UpdatePrintHandler(printSvc))
		}

		users := admin.Group("/users")
		{
			users.GET("/:id", handlers.GetUserByIDHandler(userSvc))
		}

		admin.GET("/whitelist", middleware.WhitelistEnabledMiddleware(), handlers.ListWhitelistHandler(whitelistSvc))
		admin.POST("/whitelist", middleware.WhitelistEnabledMiddleware(), handlers.AddWhitelistHandler(whitelistSvc))
		admin.DELETE("/whitelist", middleware.WhitelistEnabledMiddleware(), handlers.RemoveWhitelistHandler(whitelistSvc))
	}

	return r, nil
}
