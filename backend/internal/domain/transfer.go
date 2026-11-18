package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	TransferStatusActive  = "active"
	TransferStatusExpired = "expired"
	TransferStatusDeleted = "deleted"
	TransferStatusBlocked = "blocked"
)

type Transfer struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	ShortCode     string
	Title         *string
	Message       *string
	PasswordHash  *string
	Status        string
	ExpiresAt     time.Time
	MaxDownloads  *int32
	DownloadCount int32
	CreatedAt     time.Time
}

type File struct {
	ID           uuid.UUID
	TransferID   uuid.UUID
	OriginalName string
	StorageKey   string
	MimeType     string
	SizeBytes    int64
	CreatedAt    time.Time
}

type TransferListItem struct {
	Transfer
	FileCount      int64
	TotalSizeBytes int64
	RecipientCount int64
	FirstFileName  *string
}

type AdminTransferListItem struct {
	TransferListItem
	OwnerEmail string
}

func (t Transfer) IsExpired(now time.Time) bool {
	if t.Status == TransferStatusExpired {
		return true
	}
	return !t.ExpiresAt.After(now)
}

func (t Transfer) IsBlocked() bool {
	return t.Status == TransferStatusBlocked
}

func (t Transfer) IsDeleted() bool {
	return t.Status == TransferStatusDeleted
}

// EffectiveStatus returns the status clients should see (expiry overrides active).
func (t Transfer) EffectiveStatus(now time.Time) string {
	if t.IsDeleted() {
		return TransferStatusDeleted
	}
	if t.IsBlocked() {
		return TransferStatusBlocked
	}
	if t.IsExpired(now) {
		return TransferStatusExpired
	}
	return TransferStatusActive
}
