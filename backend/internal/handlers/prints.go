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

// NewPrintHandler godoc
// @Summary Upload a new print job STL file
// @Description Upload STL with requested filament color and create a print job
// @Tags prints
// @Accept multipart/form-data
// @Produce json
// @Param requested_filament_color formData string true "Filament color"
// @Param file formData file true "STL file to upload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security Authenticated
// @Router /prints/new [post]
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
			c.JSON(400, gin.H{"error": err.Error()})
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

// MetadataHandler godoc
// @Summary Get metadata for an STL file
// @Description Upload STL file and get its metadata (dimensions, volume, etc)
// @Tags prints
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "STL file"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security Authenticated
// @Router /metadata [post]
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
			c.JSON(500, gin.H{"error": "failed to get file metadata" + err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"metadata": metadata,
		})
	}
}

// GetUserPrintsHandler godoc
// @Summary List prints of the authenticated user
// @Description Get all print jobs for the logged-in user
// @Tags prints
// @Produce json
// @Success 200 {array} models.Print
// @Failure 401 {object} map[string]string
// @Security Authenticated
// @Router /me/prints [get]
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

// AllPrintsHandler godoc
// @Summary Get all print jobs (admin only)
// @Description List all print jobs in the system
// @Tags prints
// @Produce json
// @Success 200 {array} models.Print
// @Failure 500 {object} map[string]string
// @Security Authenticated
// @Security RoleAdmin
// @Router /prints/all [get]
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

// ApprovePrintHandler godoc
// @Summary Approve a print job (admin only)
// @Description Approves print job by ID
// @Tags prints
// @Param id path int true "Print ID"
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security Authenticated
// @Security RoleAdmin
// @Router /prints/{id}/approve [post]
func ApprovePrintHandler(printSvc *services.PrintService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		printID, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid print id"})
			return
		}

		if err := printSvc.ApprovePrint(uint(printID)); err != nil {
			c.JSON(500, gin.H{"error": "failed to approve print"})
			return
		}

		c.JSON(200, gin.H{"message": "print approved"})
	}
}

type DenyPrintRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// DenyPrintHandler godoc
// @Summary Deny a print job (admin only)
// @Description Denies print job by ID with a reason
// @Tags prints
// @Param id path int true "Print ID"
// @Accept json
// @Produce json
// @Param denyPrintRequest body handlers.DenyPrintRequest true "Denial reason"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security Authenticated
// @Security RoleAdmin
// @Router /prints/{id}/deny [post]
func DenyPrintHandler(printSvc *services.PrintService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		printID, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid print id"})
			return
		}

		var req DenyPrintRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := printSvc.DenyPrint(uint(printID), req.Reason); err != nil {
			c.JSON(500, gin.H{"error": "failed to deny print"})
			return
		}

		c.JSON(200, gin.H{"message": "print denied"})
	}
}

// DeletePrintHandler godoc
// @Summary Delete a print job and file (admin only)
// @Description Deletes print job by ID and its stored STL file
// @Tags prints
// @Param id path int true "Print ID"
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security Authenticated
// @Security RoleAdmin
// @Router /prints/{id} [delete]
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

// DownloadPrintFileHandler godoc
// @Summary Download STL file from bucket
// @Description Download STL file by filename
// @Tags prints
// @Param filename path string true "Stored filename"
// @Produce application/octet-stream
// @Success 200 "File streamed"
// @Failure 404 {object} map[string]string
// @Security Authenticated
// @Router /bucket/{filename} [get]
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

type UpdatePrintStatusRequest struct {
	Status string `json:"status" binding:"required"`
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

// UpdatePrintStatusHandler godoc
// @Summary Update the status of a print job (admin only)
// @Description Update print job status by ID (e.g. printing, completed, denied)
// @Tags prints
// @Param id path int true "Print ID"
// @Accept json
// @Produce json
// @Param updatePrintStatusRequest body handlers.UpdatePrintStatusRequest true "New status"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security Authenticated
// @Security RoleAdmin
// @Router /prints/{id}/status [put]
func UpdatePrintStatusHandler(printSvc *services.PrintService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		printID, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid print id"})
			return
		}
		var req UpdatePrintStatusRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if !isValidPrintStatus(req.Status) {
			c.JSON(400, gin.H{"error": "invalid status"})
			return
		}
		if err := printSvc.UpdatePrintStatus(uint(printID), models.PrintStatus(req.Status)); err != nil {
			c.JSON(500, gin.H{"error": "failed to update status"})
			return
		}
		c.JSON(200, gin.H{"message": "status updated"})
	}
}
