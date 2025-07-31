package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/torbenconto/spooler/internal/models"
	"github.com/torbenconto/spooler/internal/services"
	"github.com/torbenconto/spooler/internal/util"
)

type NewPrintRequest struct {
	FilamentColor string `form:"requested_filament_color" binding:"required"`
}

func NewPrintHandler(bucketSvc *services.BucketService, printSvc *services.PrintService) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}

		claims, ok := user.(*util.CustomClaims)
		if !ok {
			c.JSON(401, gin.H{"error": "invalid token claims"})
			return
		}

		var req NewPrintRequest
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request body"})
			return
		}

		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(400, gin.H{"error": "file is required"})
			return
		}

		fileHandle, err := file.Open()
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to open file"})
			return
		}
		defer fileHandle.Close()

		ext := filepath.Ext(file.Filename)
		storedFileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

		if err := bucketSvc.UploadFile(c, storedFileName, fileHandle); err != nil {
			c.JSON(500, gin.H{"error": "failed to upload file to bucket"})
			return
		}

		print := models.Print{
			UserID:                 claims.UserID,
			StoredFileName:         storedFileName,
			UploadedFileName:       file.Filename,
			RequestedFilamentColor: req.FilamentColor,
		}

		if err := printSvc.CreatePrint(&print); err != nil {
			c.JSON(500, gin.H{"error": "failed to create print"})
			return
		}

		c.JSON(200, gin.H{
			"message":                  "file uploaded successfully",
			"file":                     file.Filename,
			"backend_filename":         storedFileName,
			"requested_filament_color": req.FilamentColor,
		})
	}
}

func MetadataHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(400, gin.H{"error": "file is required"})
			return
		}

		fileHandle, err := file.Open()
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to open file"})
			return
		}
		defer fileHandle.Close()

		metadata, err := util.GetFileMetadata(file.Filename, fileHandle)
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to get file metadata"})
			return
		}

		c.JSON(200, gin.H{
			"metadata": metadata,
		})
	}
}

func GetUserPrintsHandler(printSvc *services.PrintService) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}

		claims, ok := user.(*util.CustomClaims)
		if !ok {
			c.JSON(401, gin.H{"error": "invalid token claims"})
			return
		}

		prints, err := printSvc.GetUserPrintsByID(claims.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch prints"})
			return
		}

		c.JSON(http.StatusOK, prints)
	}

}

func AllPrintsHandler(printSvc *services.PrintService) gin.HandlerFunc {
	return func(c *gin.Context) {
		prints, err := printSvc.AllPrints()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch prints"})
			return
		}

		c.JSON(http.StatusOK, prints)
	}
}

func DeletePrintHandler(printSvc *services.PrintService, bucketSvc *services.BucketService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		printID, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid print id"})
			return
		}

		printItem, err := printSvc.GetPrintByID(uint(printID))
		if err != nil {
			c.JSON(404, gin.H{"error": "print not found"})
			return
		}

		if err := bucketSvc.DeleteFile(context.Background(), printItem.StoredFileName); err != nil {
			c.JSON(500, gin.H{"error": "failed to delete file from bucket"})
			return
		}

		if err := printSvc.DeletePrint(uint(printID)); err != nil {
			c.JSON(500, gin.H{"error": "failed to delete print"})
			return
		}

		c.JSON(200, gin.H{"message": "print and file deleted"})
	}
}

func DownloadPrintFileHandler(bucketSvc *services.BucketService) gin.HandlerFunc {
	return func(c *gin.Context) {
		filename := c.Param("filename")
		reader, err := bucketSvc.GetFile(c, filename)
		if err != nil {
			c.JSON(404, gin.H{"error": "file not found"})
			return
		}
		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Header("Content-Type", "application/octet-stream")
		_, _ = io.Copy(c.Writer, reader)
	}
}

type UpdatePrintRequest struct {
	Status       string `json:"status"`
	DenialReason string `json:"denial_reason"`
}

func isValidPrintStatus(status string) bool {
	switch models.PrintStatus(status) {
	case models.StatusApprovalPending,
		models.StatusPendingPrint,
		models.StatusPrinting,
		models.StatusDenied,
		models.StatusCompleted,
		models.StatusFailed,
		models.StatusCanceled,
		models.StatusPaused:
		return true
	default:
		return false
	}
}

func UpdatePrintHandler(printSvc *services.PrintService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		printID, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid print id"})
			return
		}
		var req UpdatePrintRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request body"})
			return
		}

		updates := make(map[string]interface{})

		if req.Status != "" {
			if !isValidPrintStatus(req.Status) {
				c.JSON(400, gin.H{"error": "invalid status"})
				return
			}
			updates["status"] = req.Status
		}
		if req.DenialReason != "" {
			updates["denial_reason"] = req.DenialReason
		}

		if len(updates) == 0 {
			c.JSON(400, gin.H{"error": "no fields to update"})
			return
		}

		if err := printSvc.UpdatePrint(uint(printID), updates); err != nil {
			c.JSON(500, gin.H{"error": "failed to update print"})
			return
		}
		c.JSON(200, gin.H{"message": "print updated"})
	}
}
