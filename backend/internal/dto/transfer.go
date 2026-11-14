package dto

import (
	"time"

	"github.com/google/uuid"

	"file-temp-share/backend/internal/domain"
)

type CreateTransferRequest struct {
	Title          *string  `json:"title" binding:"omitempty,max=100"`
	Message        *string  `json:"message" binding:"omitempty,max=5000"`
	Password       *string  `json:"password" binding:"omitempty,min=4,max=72"`
	ExpiresInHours int      `json:"expires_in_hours" binding:"required,min=1,max=720"`
	Recipients     []string `json:"recipients" binding:"omitempty,max=20,dive,email"`
}

type TransferResponse struct {
	ID            uuid.UUID `json:"id"`
	ShortCode     string    `json:"short_code"`
	Title         *string   `json:"title,omitempty"`
	Message       *string   `json:"message,omitempty"`
	Status        string    `json:"status"`
	ExpiresAt     time.Time `json:"expires_at"`
	DownloadCount int32     `json:"download_count"`
	HasPassword   bool      `json:"has_password"`
	ShareURL      string    `json:"share_url"`
	CreatedAt     time.Time `json:"created_at"`
}

type TransferListItem struct {
	TransferResponse
	FileCount      int64   `json:"file_count"`
	TotalSizeBytes int64   `json:"total_size_bytes"`
	RecipientCount int64   `json:"recipient_count"`
	FirstFileName  *string `json:"first_file_name,omitempty"`
}

type TransferDetailResponse struct {
	TransferResponse
	Files          []FileResponse `json:"files"`
	FileCount      int64          `json:"file_count"`
	TotalSizeBytes int64          `json:"total_size_bytes"`
	RecipientCount int64          `json:"recipient_count"`
}

type PublicTransferResponse struct {
	ShortCode      string         `json:"short_code"`
	Title          *string        `json:"title,omitempty"`
	Message        *string        `json:"message,omitempty"`
	Status         string         `json:"status"`
	ExpiresAt      time.Time      `json:"expires_at"`
	HasPassword    bool           `json:"has_password"`
	Unlocked       bool           `json:"unlocked"`
	SenderEmail    string         `json:"sender_email"`
	Files          []FileResponse `json:"files"`
	FileCount      int64          `json:"file_count"`
	TotalSizeBytes int64          `json:"total_size_bytes"`
}

type UnlockTransferRequest struct {
	Password string `json:"password" binding:"required,min=1,max=72"`
}

type FileResponse struct {
	ID           uuid.UUID `json:"id"`
	OriginalName string    `json:"original_name"`
	MimeType     string    `json:"mime_type"`
	SizeBytes    int64     `json:"size_bytes"`
	CreatedAt    time.Time `json:"created_at"`
}

func NewTransferResponse(t domain.Transfer, shareURL string) TransferResponse {
	return TransferResponse{
		ID:            t.ID,
		ShortCode:     t.ShortCode,
		Title:         t.Title,
		Message:       t.Message,
		Status:        t.Status,
		ExpiresAt:     t.ExpiresAt,
		DownloadCount: t.DownloadCount,
		HasPassword:   t.PasswordHash != nil && *t.PasswordHash != "",
		ShareURL:      shareURL,
		CreatedAt:     t.CreatedAt,
	}
}

func NewFileResponse(f domain.File) FileResponse {
	return FileResponse{
		ID:           f.ID,
		OriginalName: f.OriginalName,
		MimeType:     f.MimeType,
		SizeBytes:    f.SizeBytes,
		CreatedAt:    f.CreatedAt,
	}
}

func NewFileResponses(files []domain.File) []FileResponse {
	out := make([]FileResponse, 0, len(files))
	for _, f := range files {
		out = append(out, NewFileResponse(f))
	}
	return out
}
