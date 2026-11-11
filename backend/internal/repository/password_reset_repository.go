package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"file-temp-share/backend/gen/sqlc"
	"file-temp-share/backend/internal/domain"
)

type PasswordResetToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}

type PasswordResetRepository struct {
	q *sqlc.Queries
}

func NewPasswordResetRepository(pool *pgxpool.Pool) *PasswordResetRepository {
	return &PasswordResetRepository{q: sqlc.New(pool)}
}

func (r *PasswordResetRepository) Create(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) (PasswordResetToken, error) {
	row, err := r.q.CreatePasswordResetToken(ctx, sqlc.CreatePasswordResetTokenParams{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return PasswordResetToken{}, fmt.Errorf("create password reset token: %w", err)
	}
	return mapResetToken(row), nil
}

func (r *PasswordResetRepository) GetByHash(ctx context.Context, tokenHash string) (PasswordResetToken, error) {
	row, err := r.q.GetPasswordResetTokenByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return PasswordResetToken{}, domain.ErrNotFound
		}
		return PasswordResetToken{}, fmt.Errorf("get password reset token: %w", err)
	}
	return mapResetToken(row), nil
}

func (r *PasswordResetRepository) MarkUsed(ctx context.Context, id uuid.UUID) error {
	if err := r.q.MarkPasswordResetTokenUsed(ctx, id); err != nil {
		return fmt.Errorf("mark password reset token used: %w", err)
	}
	return nil
}

func (r *PasswordResetRepository) InvalidateForUser(ctx context.Context, userID uuid.UUID) error {
	if err := r.q.InvalidatePasswordResetTokensForUser(ctx, userID); err != nil {
		return fmt.Errorf("invalidate password reset tokens: %w", err)
	}
	return nil
}

func mapResetToken(row sqlc.PasswordResetToken) PasswordResetToken {
	t := PasswordResetToken{
		ID:        row.ID,
		UserID:    row.UserID,
		TokenHash: row.TokenHash,
		ExpiresAt: row.ExpiresAt,
		CreatedAt: row.CreatedAt,
	}
	if row.UsedAt.Valid {
		ts := row.UsedAt.Time
		t.UsedAt = &ts
	}
	return t
}
