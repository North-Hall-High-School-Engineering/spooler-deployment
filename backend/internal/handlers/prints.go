package handlers

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

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

		const maxFileSize = 100 * 1024 * 1024
		if file.Size > maxFileSize {
			c.JSON(400, gin.H{"error": "file too large"})
			return
		}

		fileHandle, err := file.Open()
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to open file"})
			return
		}
		storedFileName := fmt.Sprintf("%s%s", uuid.New().String(), filepath.Ext(file.Filename))

		pr, pw := io.Pipe()

		go func() {
			defer fileHandle.Close()
			_, err := io.Copy(pw, fileHandle)
			_ = pw.CloseWithError(err)
		}()

		if err := bucketSvc.UploadFile(c.Request.Context(), storedFileName, pr); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload file"})
		}

		// if err := bucketSvc.UploadFile(c, storedFileName, fileHandle); err != nil {
		// 	c.JSON(500, gin.H{"error": "failed to upload file to bucket"})
		// 	return
		// }
		// ????? can i use goroutine? we will see

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

type FilePreviewType string

const (
	FileTypeSTL      FilePreviewType = "stl"
	FileType3MF      FilePreviewType = "3mf"
	FileTypeGCode3MF FilePreviewType = "gcode.3mf"
)

type FilePreview struct {
	Type FilePreviewType `json:"type"` //"stl", "3mf", "gcode.3mf"

	ModelData    *string `json:"model_data,omitempty"`    // base64, present for stl/3mf
	PreviewImage *string `json:"preview_image,omitempty"` // base64 PNG for gcode.3mf
}

func PreviewHandler() gin.HandlerFunc {
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

		var fileExtensions []string
		var fileName = file.Filename
		for {
			ext := filepath.Ext(fileName)
			if ext == "" {
				break
			}
			fileExtensions = append(fileExtensions, ext)
			fileName = strings.TrimSuffix(fileName, ext)
		}

		fileExtension := strings.Join(fileExtensions, "")
		switch fileExtension {
		case ".stl", ".3mf":
			content, err := io.ReadAll(fileHandle)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read provided file"})
				return
			}
			encoded := base64.StdEncoding.EncodeToString(content)

			var previewType FilePreviewType = FileTypeSTL
			if fileExtension == ".3mf" {
				previewType = FileType3MF
			}

			c.JSON(http.StatusOK, FilePreview{Type: previewType, ModelData: &encoded})
			return
			// case ".gcode.3mf":
		}
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
