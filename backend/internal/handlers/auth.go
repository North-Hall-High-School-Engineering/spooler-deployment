package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/torbenconto/spooler/config"
	"github.com/torbenconto/spooler/internal/models"
	"github.com/torbenconto/spooler/internal/services"
	"github.com/torbenconto/spooler/internal/util"
)

type RegisterRequest struct {
	Email     string `json:"email" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

func RegisterHandler(userSvc *services.UserService, whitelistSvc *services.WhitelistService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		// Validate email
		if !util.ValidateEmail(req.Email) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format"})
			return
		}

		if config.Cfg.EmailWhitelistEnabled {
			allowed, err := whitelistSvc.IsWhitelisted(req.Email)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "whitelist check failed"})
				return
			}
			if !allowed {
				c.JSON(http.StatusForbidden, gin.H{"error": "email not allowed"})
				return
			}
		}

		user := models.User{
			Email:     req.Email,
			FirstName: req.FirstName,
			LastName:  req.LastName,
		}

		if err := userSvc.CreateUser(&user); err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "user registered"})
	}
}

func MeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "authenticated", "user": user})
	}
}

func GetUserByIDHandler(userSvc *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		userID, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid print id"})
			return
		}

		user, err := userSvc.GetUserByID(uint(userID))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		}

		c.JSON(http.StatusOK, gin.H{"message": "ok", "user": user})
	}
}
