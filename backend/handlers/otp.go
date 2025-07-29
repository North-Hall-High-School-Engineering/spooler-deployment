package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/torbenconto/spooler/config"
	"github.com/torbenconto/spooler/models"
	"github.com/torbenconto/spooler/services"
	"github.com/torbenconto/spooler/util"
)

type RequestOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func RequestOTPHandler(otpSvc *services.OTPService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RequestOTPRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		code, err := otpSvc.GenerateCode(req.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate code"})
			return
		}

		go func() {
			err := (&util.EmailSender{
				From:     config.Cfg.SMTPEmail,
				Password: config.Cfg.SMTPPassword,
				SMTPHost: config.Cfg.SMTPHost,
				SMTPPort: config.Cfg.SMTPPort,
			}).Send(util.EmailMessage{
				To:      req.Email,
				Subject: "SP00LER: One Time Passcode",
				Body:    fmt.Sprintf("Your one time passcode is %s", code),
			})
			if err != nil {
				log.Printf("error sending OTP to %s: %v", req.Email, err)
			}
		}()

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
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "token",
			Value:    token,
			MaxAge:   3600 * 24 * 7,
			Secure:   secure,
			Path:     "/",
			HttpOnly: true,
		})

		c.JSON(http.StatusOK, gin.H{"message": "login successful", "token": token})
	}
}
