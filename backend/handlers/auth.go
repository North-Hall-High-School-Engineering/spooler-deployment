package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/torbenconto/spooler/models"
	"github.com/torbenconto/spooler/services"
	"github.com/torbenconto/spooler/util"
)

type RegisterRequest struct {
	Email     string `json:"email" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

func RegisterHandler(userSvc *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate email
		if !util.ValidateEmail(req.Email) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format"})
		}

		user := models.User{
			Email:     req.Email,
			FirstName: req.FirstName,
			LastName:  req.LastName,
		}

		if err := userSvc.CreateUser(&user); err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "user registered"})
	}
}

func MeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}

		c.JSON(200, gin.H{"message": "authenticated", "user": user})
	}
}
