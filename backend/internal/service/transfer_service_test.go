package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"file-temp-share/backend/internal/auth"
	"file-temp-share/backend/internal/domain"
)

func TestEnsureDownloadable(t *testing.T) {
	svc := &TransferService{}
	now := time.Now().UTC()

	t.Run("blocked", func(t *testing.T) {
		err := svc.ensureDownloadable(domain.Transfer{
			Status:    domain.TransferStatusBlocked,
			ExpiresAt: now.Add(time.Hour),
		})
		if !errors.Is(err, domain.ErrTransferBlocked) {
			t.Fatalf("got %v, want ErrTransferBlocked", err)
		}
	})

	t.Run("expired by time", func(t *testing.T) {
		err := svc.ensureDownloadable(domain.Transfer{
			Status:    domain.TransferStatusActive,
			ExpiresAt: now.Add(-time.Minute),
		})
		if !errors.Is(err, domain.ErrTransferExpired) {
			t.Fatalf("got %v, want ErrTransferExpired", err)
		}
	})

	t.Run("download limit reached", func(t *testing.T) {
		max := int32(2)
		err := svc.ensureDownloadable(domain.Transfer{
			Status:        domain.TransferStatusActive,
			ExpiresAt:     now.Add(time.Hour),
			MaxDownloads:  &max,
			DownloadCount: 2,
		})
		if !errors.Is(err, domain.ErrDownloadLimitReached) {
			t.Fatalf("got %v, want ErrDownloadLimitReached", err)
		}
	})

	t.Run("ok", func(t *testing.T) {
		max := int32(5)
		err := svc.ensureDownloadable(domain.Transfer{
			Status:        domain.TransferStatusActive,
			ExpiresAt:     now.Add(time.Hour),
			MaxDownloads:  &max,
			DownloadCount: 1,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestCheckTransferPassword(t *testing.T) {
	svc := &TransferService{}
	hash, err := auth.HashPassword("secret")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}

	t.Run("no password required", func(t *testing.T) {
		if err := svc.checkTransferPassword(domain.Transfer{}, ""); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("password required", func(t *testing.T) {
		err := svc.checkTransferPassword(domain.Transfer{PasswordHash: &hash}, "  ")
		if !errors.Is(err, domain.ErrTransferPasswordRequired) {
			t.Fatalf("got %v, want ErrTransferPasswordRequired", err)
		}
	})

	t.Run("invalid password", func(t *testing.T) {
		err := svc.checkTransferPassword(domain.Transfer{PasswordHash: &hash}, "wrong")
		if !errors.Is(err, domain.ErrInvalidTransferPassword) {
			t.Fatalf("got %v, want ErrInvalidTransferPassword", err)
		}
	})

	t.Run("valid password", func(t *testing.T) {
		if err := svc.checkTransferPassword(domain.Transfer{PasswordHash: &hash}, "secret"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestUploadFileSizeRules(t *testing.T) {
	svc := &TransferService{uploadMaxBytes: 100}

	t.Run("empty file", func(t *testing.T) {
		_, err := svc.UploadFile(context.Background(), UploadFileInput{
			TransferID: uuid.New(),
			UserID:     uuid.New(),
			Size:       0,
			Body:       strings.NewReader(""),
		})
		if !errors.Is(err, domain.ErrEmptyFile) {
			t.Fatalf("got %v, want ErrEmptyFile", err)
		}
	})

	t.Run("too large", func(t *testing.T) {
		_, err := svc.UploadFile(context.Background(), UploadFileInput{
			TransferID: uuid.New(),
			UserID:     uuid.New(),
			Size:       101,
			Body:       strings.NewReader(strings.Repeat("a", 101)),
		})
		if !errors.Is(err, domain.ErrFileTooLarge) {
			t.Fatalf("got %v, want ErrFileTooLarge", err)
		}
	})
}
