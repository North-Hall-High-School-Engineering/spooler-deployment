package main

import (
	"context"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/torbenconto/spooler/config"
	"github.com/torbenconto/spooler/internal/handlers"
	"github.com/torbenconto/spooler/internal/middleware"
	"github.com/torbenconto/spooler/internal/models"
	"github.com/torbenconto/spooler/internal/services"
	"github.com/torbenconto/spooler/internal/storage"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB) (*gin.Engine, error) {
	r := gin.Default()

	var allowedOrigins []string
	if len(config.Cfg.CORSAllowOrigins) == 0 {
		allowedOrigins = []string{"http://localhost:5173"}
	}

	for _, origin := range config.Cfg.CORSAllowOrigins {
		allowedOrigins = append(allowedOrigins, strings.TrimSpace(origin))
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length", "Set-Cookie"},
		MaxAge:           12 * time.Hour,
	}))

	storageClient, err := storage.NewStorageClient(context.Background(), config.Cfg)
	if err != nil {
		return nil, err
	}

	userSvc := services.NewUserService(db)
	otpSvc := services.NewOTPService(db)
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
		auth.GET("/bucket/:filename", handlers.DownloadPrintFileHandler(storageClient))
		auth.POST("/preview", handlers.PreviewHandler())
		auth.POST("/prints/new", handlers.NewPrintHandler(storageClient, printSvc))
	}

	// Admin-only routes
	admin := auth.Group("/")
	admin.Use(middleware.RoleAuthMiddleware(models.RoleAdmin))
	{
		prints := admin.Group("/prints")
		{
			prints.GET("/all", handlers.AllPrintsHandler(printSvc))

			prints.DELETE("/:id", handlers.DeletePrintHandler(printSvc, storageClient))
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
