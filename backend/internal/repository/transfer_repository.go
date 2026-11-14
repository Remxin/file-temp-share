package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"file-temp-share/backend/gen/sqlc"
	"file-temp-share/backend/internal/domain"
)

type TransferRepository struct {
	q *sqlc.Queries
}

func NewTransferRepository(pool *pgxpool.Pool) *TransferRepository {
	return &TransferRepository{q: sqlc.New(pool)}
}

type CreateTransferParams struct {
	UserID       uuid.UUID
	ShortCode    string
	Title        *string
	Message      *string
	PasswordHash *string
	ExpiresAt    time.Time
}

func (r *TransferRepository) Create(ctx context.Context, p CreateTransferParams) (domain.Transfer, error) {
	row, err := r.q.CreateTransfer(ctx, sqlc.CreateTransferParams{
		UserID:       p.UserID,
		ShortCode:    p.ShortCode,
		Title:        textPtr(p.Title),
		Message:      textPtr(p.Message),
		PasswordHash: textPtr(p.PasswordHash),
		ExpiresAt:    p.ExpiresAt,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.Transfer{}, domain.ErrConflict
		}
		return domain.Transfer{}, fmt.Errorf("create transfer: %w", err)
	}
	return mapTransfer(row), nil
}

func (r *TransferRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Transfer, error) {
	row, err := r.q.GetTransferByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Transfer{}, domain.ErrNotFound
		}
		return domain.Transfer{}, fmt.Errorf("get transfer by id: %w", err)
	}
	return mapTransfer(row), nil
}

func (r *TransferRepository) GetByShortCode(ctx context.Context, shortCode string) (domain.Transfer, error) {
	row, err := r.q.GetTransferByShortCode(ctx, shortCode)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Transfer{}, domain.ErrNotFound
		}
		return domain.Transfer{}, fmt.Errorf("get transfer by short code: %w", err)
	}
	return mapTransfer(row), nil
}

func (r *TransferRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]domain.TransferListItem, error) {
	rows, err := r.q.ListTransfersWithStatsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list transfers: %w", err)
	}
	out := make([]domain.TransferListItem, 0, len(rows))
	for _, row := range rows {
		item := domain.TransferListItem{
			Transfer:       mapTransferFromListRow(row),
			FileCount:      row.FileCount,
			TotalSizeBytes: row.TotalSizeBytes,
			RecipientCount: row.RecipientCount,
		}
		if row.FirstFileName != "" {
			name := row.FirstFileName
			item.FirstFileName = &name
		}
		out = append(out, item)
	}
	return out, nil
}

func (r *TransferRepository) ListFiles(ctx context.Context, transferID uuid.UUID) ([]domain.File, error) {
	rows, err := r.q.ListFilesByTransferID(ctx, transferID)
	if err != nil {
		return nil, fmt.Errorf("list files: %w", err)
	}
	out := make([]domain.File, 0, len(rows))
	for _, row := range rows {
		out = append(out, mapFile(row))
	}
	return out, nil
}

func (r *TransferRepository) CountRecipients(ctx context.Context, transferID uuid.UUID) (int64, error) {
	n, err := r.q.CountRecipientsByTransferID(ctx, transferID)
	if err != nil {
		return 0, fmt.Errorf("count recipients: %w", err)
	}
	return n, nil
}

func (r *TransferRepository) AddRecipient(ctx context.Context, transferID uuid.UUID, email string) error {
	_, err := r.q.CreateTransferRecipient(ctx, sqlc.CreateTransferRecipientParams{
		TransferID: transferID,
		Email:      email,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil // duplicate recipient is fine
		}
		return fmt.Errorf("create transfer recipient: %w", err)
	}
	return nil
}

func (r *TransferRepository) CreateFile(ctx context.Context, transferID uuid.UUID, originalName, storageKey, mimeType string, sizeBytes int64) (domain.File, error) {
	row, err := r.q.CreateFile(ctx, sqlc.CreateFileParams{
		TransferID:   transferID,
		OriginalName: originalName,
		StorageKey:   storageKey,
		MimeType:     mimeType,
		SizeBytes:    sizeBytes,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.File{}, domain.ErrConflict
		}
		return domain.File{}, fmt.Errorf("create file: %w", err)
	}
	return mapFile(row), nil
}

func mapTransferFromListRow(row sqlc.ListTransfersWithStatsByUserIDRow) domain.Transfer {
	t := domain.Transfer{
		ID:            row.ID,
		UserID:        row.UserID,
		ShortCode:     row.ShortCode,
		Title:         textToPtr(row.Title),
		Message:       textToPtr(row.Message),
		PasswordHash:  textToPtr(row.PasswordHash),
		Status:        row.Status,
		ExpiresAt:     row.ExpiresAt,
		DownloadCount: row.DownloadCount,
		CreatedAt:     row.CreatedAt,
	}
	if row.MaxDownloads.Valid {
		v := row.MaxDownloads.Int32
		t.MaxDownloads = &v
	}
	return t
}

func mapTransfer(row sqlc.Transfer) domain.Transfer {
	t := domain.Transfer{
		ID:            row.ID,
		UserID:        row.UserID,
		ShortCode:     row.ShortCode,
		Title:         textToPtr(row.Title),
		Message:       textToPtr(row.Message),
		PasswordHash:  textToPtr(row.PasswordHash),
		Status:        row.Status,
		ExpiresAt:     row.ExpiresAt,
		DownloadCount: row.DownloadCount,
		CreatedAt:     row.CreatedAt,
	}
	if row.MaxDownloads.Valid {
		v := row.MaxDownloads.Int32
		t.MaxDownloads = &v
	}
	return t
}

func mapFile(row sqlc.File) domain.File {
	return domain.File{
		ID:           row.ID,
		TransferID:   row.TransferID,
		OriginalName: row.OriginalName,
		StorageKey:   row.StorageKey,
		MimeType:     row.MimeType,
		SizeBytes:    row.SizeBytes,
		CreatedAt:    row.CreatedAt,
	}
}

func textPtr(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func textToPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	s := t.String
	return &s
}
