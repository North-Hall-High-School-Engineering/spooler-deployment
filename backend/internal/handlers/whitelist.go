package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/torbenconto/spooler/internal/services"
)

func ListWhitelistHandler(svc *services.WhitelistService) gin.HandlerFunc {
	return func(c *gin.Context) {
		list, err := svc.List()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch whitelist"})
			return
		}
		c.JSON(http.StatusOK, list)
	}
}

type AddWhiteListRequest struct {
	Emails []string `json:"emails"`
}

func AddWhitelistHandler(whitelistSvc *services.WhitelistService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AddWhiteListRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		if len(req.Emails) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no emails provided"})
			return
		}
		if err := whitelistSvc.Add(req.Emails...); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add emails to whitelist"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "added"})
	}
}

type RemoveWhitelistRequest struct {
	Email string `json:"email"`
}

func RemoveWhitelistHandler(svc *services.WhitelistService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RemoveWhitelistRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		if err := svc.Remove(req.Email); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not remove"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "removed"})
	}
}
