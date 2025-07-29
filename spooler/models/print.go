package models

import "time"

type PrintStatus string

const (
	StatusApprovalPending PrintStatus = "approval_pending"
	StatusPendingPrint    PrintStatus = "pending_print"
	StatusPrinting        PrintStatus = "printing"
	StatusDenied          PrintStatus = "denied"
	StatusCompleted       PrintStatus = "completed"
	StatusFailed          PrintStatus = "failed"
	StatusCanceled        PrintStatus = "canceled"
	StatusPaused          PrintStatus = "paused"
)

type Print struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint `gorm:"index"`

	Status                 PrintStatus `gorm:"type:varchar(32);default:'approval_pending'"`
	Progress               int         `gorm:"default:0"`
	UploadedFileName       string      `gorm:"not null"`
	StoredFileName         string      `gorm:"not null"`
	RequestedFilamentColor string      `gorm:"not null;default:'#000000'"`
	DenialReason           string

	CreatedAt time.Time
	UpdatedAt time.Time
}
