package dto

import (
	"github.com/google/uuid"

	"file-temp-share/backend/internal/domain"
)

type AdminSetUserActiveRequest struct {
	IsActive *bool `json:"is_active" binding:"required"`
}

type AdminSetTransferStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=blocked active"`
}

type AdminTransferListItem struct {
	TransferListItem
	UserID     uuid.UUID `json:"user_id"`
	OwnerEmail string    `json:"owner_email"`
}

func NewAdminTransferListItem(item domain.AdminTransferListItem, shareURL string) AdminTransferListItem {
	base := NewTransferResponse(item.Transfer, shareURL)
	return AdminTransferListItem{
		TransferListItem: TransferListItem{
			TransferResponse: base,
			FileCount:        item.FileCount,
			TotalSizeBytes:   item.TotalSizeBytes,
			RecipientCount:   item.RecipientCount,
			FirstFileName:    item.FirstFileName,
		},
		UserID:     item.UserID,
		OwnerEmail: item.OwnerEmail,
	}
}
