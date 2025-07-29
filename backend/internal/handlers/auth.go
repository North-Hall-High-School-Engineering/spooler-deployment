package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/torbenconto/spooler/internal/models"
	"github.com/torbenconto/spooler/internal/services"
	"github.com/torbenconto/spooler/internal/util"
)

type RegisterRequest struct {
	Email     string `json:"email" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// RegisterHandler godoc
// @Summary Register a new user account
// @Description Creates a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param registerRequest body handlers.RegisterRequest true "User registration payload"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /register [post]
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

// MeHandler godoc
// @Summary Get authenticated user info
// @Description Returns current user's info from JWT claims
// @Tags user
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Security Authenticated
// @Router /me [get]
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
