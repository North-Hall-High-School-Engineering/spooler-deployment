package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/torbenconto/spooler/config"
	"github.com/torbenconto/spooler/internal/models"
	"github.com/torbenconto/spooler/internal/services"
	"github.com/torbenconto/spooler/internal/util"
)

type RequestOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func RequestOTPHandler(otpSvc *services.OTPService, whitelistSvc *services.WhitelistService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RequestOTPRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		if config.Cfg.Features.EmailWhitelistEnabled {
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

		code, err := otpSvc.GenerateCode(req.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate code"})
			return
		}

		err = (&util.EmailSender{
			From:     config.Cfg.SMTP.Email,
			Password: config.Cfg.SMTP.Password,
			SMTPHost: config.Cfg.SMTP.Host,
			SMTPPort: config.Cfg.SMTP.Port,
		}).Send(util.EmailMessage{
			To:      req.Email,
			Subject: "SP00LER: One Time Passcode",
			Body:    fmt.Sprintf("Your one time passcode is %s", code),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send OTP"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "OTP sent to email"})
	}
}

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

func VerifyOTPHandler(otpSvc *services.OTPService, userSvc *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req VerifyOTPRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		if err := otpSvc.ValidateCode(req.Email, req.Code); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired code"})
			return
		}

		user, err := userSvc.GetUserByEmail(req.Email)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "no user found under given email"})
			return
		}

		token, err := util.GenerateJWT(req.Email, models.Role(user.Role), user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate JWT token"})
		}

		var secure = gin.Mode() == gin.ReleaseMode
		cookie := &http.Cookie{
			Name:     "token",
			Value:    token,
			MaxAge:   3600 * 24 * 7,
			Path:     "/",
			HttpOnly: true,
			Secure:   secure,
		}

		if secure {
			cookie.SameSite = http.SameSiteNoneMode
		}
		http.SetCookie(c.Writer, cookie)

		c.JSON(http.StatusOK, gin.H{"message": "login successful", "token": token})
	}
}
