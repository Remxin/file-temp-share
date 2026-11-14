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
	"file-temp-share/backend/internal/queue"
	"file-temp-share/backend/internal/repository"
	"file-temp-share/backend/internal/storage"
)

const (
	shortCodeLen      = 10
	shortCodeAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	createRetries     = 5
)

type TransferService struct {
	transfers      *repository.TransferRepository
	users          *repository.UserRepository
	blobs          storage.Storage
	publisher      queue.Publisher
	appPublicURL   string
	uploadMaxBytes int64
}

func NewTransferService(
	transfers *repository.TransferRepository,
	users *repository.UserRepository,
	blobs storage.Storage,
	publisher queue.Publisher,
	appPublicURL string,
	uploadMaxBytes int64,
) *TransferService {
	return &TransferService{
		transfers:      transfers,
		users:          users,
		blobs:          blobs,
		publisher:      publisher,
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

type TransferDetail struct {
	Transfer       domain.Transfer
	Files          []domain.File
	FileCount      int64
	TotalSizeBytes int64
	RecipientCount int64
}

type PublicTransfer struct {
	Transfer       domain.Transfer
	Files          []domain.File
	FileCount      int64
	TotalSizeBytes int64
	SenderEmail    string
	Unlocked       bool
}

type DownloadResult struct {
	File     domain.File
	Body     io.ReadCloser
	MIMEType string
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

	// Notify recipients once the first file makes the transfer downloadable.
	files, listErr := s.transfers.ListFiles(ctx, transfer.ID)
	if listErr != nil {
		slog.Warn("list files after upload failed", "transfer_id", transfer.ID, "err", listErr)
	} else if len(files) == 1 {
		s.enqueueRecipientInvites(ctx, transfer, file, len(files))
	}

	return file, nil
}

func (s *TransferService) enqueueRecipientInvites(
	ctx context.Context,
	transfer domain.Transfer,
	file domain.File,
	fileCount int,
) {
	if s.publisher == nil {
		return
	}
	emails, err := s.transfers.ListRecipientEmails(ctx, transfer.ID)
	if err != nil {
		slog.Warn("list recipients for invite failed", "transfer_id", transfer.ID, "err", err)
		return
	}
	if len(emails) == 0 {
		return
	}

	title := ""
	if transfer.Title != nil {
		title = *transfer.Title
	}
	message := ""
	if transfer.Message != nil {
		message = *transfer.Message
	}
	downloadURL := s.ShareURL(transfer)
	sentAt := time.Now().UTC()

	for _, email := range emails {
		msg := queue.TransferInviteEmailMessage{
			Type:          queue.TypeTransferInvite,
			To:            email,
			DownloadURL:   downloadURL,
			Title:         title,
			Message:       message,
			FileName:      file.OriginalName,
			FileSizeBytes: file.SizeBytes,
			FileCount:     fileCount,
			ExpiresAt:     transfer.ExpiresAt,
			SentAt:        sentAt,
		}
		if err := s.publisher.PublishTransferInvite(ctx, msg); err != nil {
			slog.Warn("enqueue transfer invite failed",
				"transfer_id", transfer.ID,
				"to", email,
				"err", err,
			)
			continue
		}
		slog.Info("transfer invite email enqueued",
			"transfer_id", transfer.ID,
			"to", email,
		)
	}
}

func (s *TransferService) List(ctx context.Context, userID uuid.UUID, filter string) ([]domain.TransferListItem, error) {
	items, err := s.transfers.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	out := make([]domain.TransferListItem, 0, len(items))
	for _, item := range items {
		status := item.EffectiveStatus(now)
		item.Status = status
		switch filter {
		case "active":
			if status != domain.TransferStatusActive {
				continue
			}
		case "expired":
			if status != domain.TransferStatusExpired {
				continue
			}
		case "all", "":
			// keep
		default:
			return nil, fmt.Errorf("%w: invalid status filter", domain.ErrValidation)
		}
		out = append(out, item)
	}
	return out, nil
}

func (s *TransferService) GetOwned(ctx context.Context, userID, transferID uuid.UUID) (TransferDetail, error) {
	transfer, err := s.transfers.GetByID(ctx, transferID)
	if err != nil {
		return TransferDetail{}, err
	}
	if transfer.IsDeleted() {
		return TransferDetail{}, domain.ErrNotFound
	}
	if transfer.UserID != userID {
		return TransferDetail{}, domain.ErrForbidden
	}
	transfer.Status = transfer.EffectiveStatus(time.Now().UTC())

	files, err := s.transfers.ListFiles(ctx, transfer.ID)
	if err != nil {
		return TransferDetail{}, err
	}
	var total int64
	for _, f := range files {
		total += f.SizeBytes
	}
	recipients, err := s.transfers.CountRecipients(ctx, transfer.ID)
	if err != nil {
		return TransferDetail{}, err
	}
	return TransferDetail{
		Transfer:       transfer,
		Files:          files,
		FileCount:      int64(len(files)),
		TotalSizeBytes: total,
		RecipientCount: recipients,
	}, nil
}

func (s *TransferService) GetPublic(ctx context.Context, shortCode string) (PublicTransfer, error) {
	return s.loadPublic(ctx, shortCode, "", false)
}

func (s *TransferService) UnlockPublic(ctx context.Context, shortCode, password string) (PublicTransfer, error) {
	return s.loadPublic(ctx, shortCode, password, true)
}

func (s *TransferService) loadPublic(ctx context.Context, shortCode, password string, requireUnlock bool) (PublicTransfer, error) {
	transfer, err := s.transfers.GetByShortCode(ctx, shortCode)
	if err != nil {
		return PublicTransfer{}, err
	}
	if transfer.IsDeleted() {
		return PublicTransfer{}, domain.ErrNotFound
	}
	if transfer.IsBlocked() {
		return PublicTransfer{}, domain.ErrTransferBlocked
	}
	transfer.Status = transfer.EffectiveStatus(time.Now().UTC())

	hasPassword := transfer.PasswordHash != nil && *transfer.PasswordHash != ""
	unlocked := !hasPassword
	if hasPassword && password != "" {
		if !auth.CheckPassword(*transfer.PasswordHash, password) {
			return PublicTransfer{}, domain.ErrInvalidTransferPassword
		}
		unlocked = true
	}
	if requireUnlock && hasPassword && !unlocked {
		return PublicTransfer{}, domain.ErrTransferPasswordRequired
	}

	senderEmail := ""
	if user, userErr := s.users.GetByID(ctx, transfer.UserID); userErr == nil {
		senderEmail = user.Email
	}

	out := PublicTransfer{
		Transfer:    transfer,
		SenderEmail: senderEmail,
		Unlocked:    unlocked,
		Files:       []domain.File{},
	}

	files, err := s.transfers.ListFiles(ctx, transfer.ID)
	if err != nil {
		return PublicTransfer{}, err
	}
	var total int64
	for _, f := range files {
		total += f.SizeBytes
	}
	out.FileCount = int64(len(files))
	out.TotalSizeBytes = total

	// Password-protected: hide file list until unlocked.
	if unlocked {
		out.Files = files
	}
	return out, nil
}

func (s *TransferService) DownloadPublicFile(
	ctx context.Context,
	shortCode string,
	fileID uuid.UUID,
	password string,
	ip, userAgent string,
) (DownloadResult, error) {
	transfer, err := s.transfers.GetByShortCode(ctx, shortCode)
	if err != nil {
		return DownloadResult{}, err
	}
	if err := s.ensureDownloadable(transfer); err != nil {
		return DownloadResult{}, err
	}
	if err := s.checkTransferPassword(transfer, password); err != nil {
		return DownloadResult{}, err
	}
	return s.openAndLogDownload(ctx, transfer, fileID, ip, userAgent)
}

func (s *TransferService) DownloadOwnedFile(
	ctx context.Context,
	userID, transferID, fileID uuid.UUID,
	ip, userAgent string,
) (DownloadResult, error) {
	transfer, err := s.transfers.GetByID(ctx, transferID)
	if err != nil {
		return DownloadResult{}, err
	}
	if transfer.IsDeleted() {
		return DownloadResult{}, domain.ErrNotFound
	}
	if transfer.UserID != userID {
		return DownloadResult{}, domain.ErrForbidden
	}
	if err := s.ensureDownloadable(transfer); err != nil {
		return DownloadResult{}, err
	}
	return s.openAndLogDownload(ctx, transfer, fileID, ip, userAgent)
}

func (s *TransferService) ensureDownloadable(transfer domain.Transfer) error {
	if transfer.IsBlocked() {
		return domain.ErrTransferBlocked
	}
	if transfer.IsExpired(time.Now().UTC()) {
		return domain.ErrTransferExpired
	}
	if transfer.MaxDownloads != nil && transfer.DownloadCount >= *transfer.MaxDownloads {
		return domain.ErrDownloadLimitReached
	}
	return nil
}

func (s *TransferService) checkTransferPassword(transfer domain.Transfer, password string) error {
	if transfer.PasswordHash == nil || *transfer.PasswordHash == "" {
		return nil
	}
	if strings.TrimSpace(password) == "" {
		return domain.ErrTransferPasswordRequired
	}
	if !auth.CheckPassword(*transfer.PasswordHash, password) {
		return domain.ErrInvalidTransferPassword
	}
	return nil
}

func (s *TransferService) openAndLogDownload(
	ctx context.Context,
	transfer domain.Transfer,
	fileID uuid.UUID,
	ip, userAgent string,
) (DownloadResult, error) {
	file, err := s.transfers.GetFileByID(ctx, fileID)
	if err != nil {
		return DownloadResult{}, err
	}
	if file.TransferID != transfer.ID {
		return DownloadResult{}, domain.ErrNotFound
	}

	body, err := s.blobs.Get(ctx, file.StorageKey)
	if err != nil {
		return DownloadResult{}, fmt.Errorf("storage get: %w", err)
	}

	if _, err := s.transfers.IncrementDownloadCount(ctx, transfer.ID); err != nil {
		_ = body.Close()
		return DownloadResult{}, err
	}
	if err := s.transfers.CreateDownloadLog(ctx, transfer.ID, &file.ID, ip, userAgent); err != nil {
		slog.Warn("download log failed", "transfer_id", transfer.ID, "file_id", file.ID, "err", err)
	}

	return DownloadResult{
		File:     file,
		Body:     body,
		MIMEType: file.MimeType,
	}, nil
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
