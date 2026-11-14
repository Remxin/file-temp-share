package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"file-temp-share/backend/internal/auth"
	"file-temp-share/backend/internal/domain"
	"file-temp-share/backend/internal/repository"
	"file-temp-share/backend/internal/storage"
)

const (
	shortCodeLen     = 10
	shortCodeAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	createRetries    = 5
)

type TransferService struct {
	transfers      *repository.TransferRepository
	blobs          storage.Storage
	appPublicURL   string
	uploadMaxBytes int64
}

func NewTransferService(
	transfers *repository.TransferRepository,
	blobs storage.Storage,
	appPublicURL string,
	uploadMaxBytes int64,
) *TransferService {
	return &TransferService{
		transfers:      transfers,
		blobs:          blobs,
		appPublicURL:   strings.TrimRight(appPublicURL, "/"),
		uploadMaxBytes: uploadMaxBytes,
	}
}

type CreateTransferInput struct {
	Title          *string
	Message        *string
	Password       *string
	ExpiresInHours int
	Recipients     []string
}

func (s *TransferService) Create(ctx context.Context, userID uuid.UUID, in CreateTransferInput) (domain.Transfer, error) {
	var passwordHash *string
	if in.Password != nil && strings.TrimSpace(*in.Password) != "" {
		hash, err := auth.HashPassword(strings.TrimSpace(*in.Password))
		if err != nil {
			return domain.Transfer{}, err
		}
		passwordHash = &hash
	}

	expiresAt := time.Now().UTC().Add(time.Duration(in.ExpiresInHours) * time.Hour)

	var transfer domain.Transfer
	var err error
	for i := 0; i < createRetries; i++ {
		code, genErr := randomShortCode(shortCodeLen)
		if genErr != nil {
			return domain.Transfer{}, genErr
		}
		transfer, err = s.transfers.Create(ctx, repository.CreateTransferParams{
			UserID:       userID,
			ShortCode:    code,
			Title:        trimPtr(in.Title),
			Message:      trimPtr(in.Message),
			PasswordHash: passwordHash,
			ExpiresAt:    expiresAt,
		})
		if err == nil {
			break
		}
		if errors.Is(err, domain.ErrConflict) {
			continue
		}
		return domain.Transfer{}, err
	}
	if err != nil {
		return domain.Transfer{}, fmt.Errorf("create transfer: exhausted short_code retries: %w", err)
	}

	for _, email := range in.Recipients {
		email = strings.ToLower(strings.TrimSpace(email))
		if email == "" {
			continue
		}
		if addErr := s.transfers.AddRecipient(ctx, transfer.ID, email); addErr != nil {
			slog.Warn("add transfer recipient failed", "transfer_id", transfer.ID, "err", addErr)
		}
	}

	return transfer, nil
}

type UploadFileInput struct {
	TransferID   uuid.UUID
	UserID       uuid.UUID
	OriginalName string
	ContentType  string
	Size         int64
	Body         io.Reader
}

func (s *TransferService) UploadFile(ctx context.Context, in UploadFileInput) (domain.File, error) {
	if in.Size <= 0 {
		return domain.File{}, domain.ErrEmptyFile
	}
	if in.Size > s.uploadMaxBytes {
		return domain.File{}, domain.ErrFileTooLarge
	}

	transfer, err := s.transfers.GetByID(ctx, in.TransferID)
	if err != nil {
		return domain.File{}, err
	}
	if transfer.UserID != in.UserID {
		return domain.File{}, domain.ErrForbidden
	}
	if transfer.IsDeleted() {
		return domain.File{}, domain.ErrNotFound
	}
	if transfer.IsBlocked() {
		return domain.File{}, domain.ErrTransferBlocked
	}
	if transfer.IsExpired(time.Now().UTC()) {
		return domain.File{}, domain.ErrTransferExpired
	}

	name := sanitizeFilename(in.OriginalName)
	if name == "" {
		name = "file"
	}
	contentType := strings.TrimSpace(in.ContentType)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	storageKey := fmt.Sprintf("transfers/%s/%s", transfer.ID.String(), uuid.NewString())
	if err := s.blobs.Put(ctx, storageKey, in.Body, in.Size, contentType); err != nil {
		return domain.File{}, fmt.Errorf("storage put: %w", err)
	}

	file, err := s.transfers.CreateFile(ctx, transfer.ID, name, storageKey, contentType, in.Size)
	if err != nil {
		if delErr := s.blobs.Delete(ctx, storageKey); delErr != nil {
			slog.Error("orphan object cleanup failed", "key", storageKey, "err", delErr)
		}
		return domain.File{}, err
	}
	return file, nil
}

func (s *TransferService) ShareURL(t domain.Transfer) string {
	return s.appPublicURL + "/t/" + t.ShortCode
}

func randomShortCode(n int) (string, error) {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("short code: %w", err)
	}
	out := make([]byte, n)
	for i := 0; i < n; i++ {
		out[i] = shortCodeAlphabet[int(buf[i])%len(shortCodeAlphabet)]
	}
	return string(out), nil
}

func trimPtr(s *string) *string {
	if s == nil {
		return nil
	}
	v := strings.TrimSpace(*s)
	if v == "" {
		return nil
	}
	return &v
}

func sanitizeFilename(name string) string {
	name = filepath.Base(strings.TrimSpace(name))
	name = strings.ReplaceAll(name, "\x00", "")
	if name == "." || name == ".." {
		return ""
	}
	return name
}
