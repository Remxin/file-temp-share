package service

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"

	"file-temp-share/backend/internal/auth"
	"file-temp-share/backend/internal/domain"
	"file-temp-share/backend/internal/repository"
)

const passwordResetTTL = time.Hour

type AuthService struct {
	users    *repository.UserRepository
	resets   *repository.PasswordResetRepository
	tokens   *auth.TokenManager
	resetTTL time.Duration
}

func NewAuthService(
	users *repository.UserRepository,
	resets *repository.PasswordResetRepository,
	tokens *auth.TokenManager,
) *AuthService {
	return &AuthService{
		users:    users,
		resets:   resets,
		tokens:   tokens,
		resetTTL: passwordResetTTL,
	}
}

type AuthResult struct {
	User      domain.User
	Token     string
	ExpiresAt time.Time
}

func (s *AuthService) Register(ctx context.Context, email, password string) (AuthResult, error) {
	email = normalizeEmail(email)
	hash, err := auth.HashPassword(password)
	if err != nil {
		return AuthResult{}, err
	}
	user, err := s.users.Create(ctx, email, hash, domain.RoleUser)
	if err != nil {
		return AuthResult{}, err
	}
	return s.issue(user)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (AuthResult, error) {
	email = normalizeEmail(email)
	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return AuthResult{}, domain.ErrInvalidCredentials
		}
		return AuthResult{}, err
	}
	if !user.IsActive {
		return AuthResult{}, domain.ErrUserInactive
	}
	if !auth.CheckPassword(user.PasswordHash, password) {
		return AuthResult{}, domain.ErrInvalidCredentials
	}
	return s.issue(user)
}

func (s *AuthService) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	user, err := s.users.GetByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	if !user.IsActive {
		return domain.User{}, domain.ErrUserInactive
	}
	return user, nil
}

// RequestPasswordReset always succeeds from the caller's perspective (no email enumeration).
// Mail sending is not implemented yet — the raw token is logged for local/dev use.
func (s *AuthService) RequestPasswordReset(ctx context.Context, email string) error {
	email = normalizeEmail(email)
	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil
		}
		return err
	}
	if !user.IsActive {
		return nil
	}

	raw, hash, err := auth.NewResetToken()
	if err != nil {
		return err
	}
	expiresAt := time.Now().UTC().Add(s.resetTTL)
	if err := s.resets.InvalidateForUser(ctx, user.ID); err != nil {
		return err
	}
	if _, err := s.resets.Create(ctx, user.ID, hash, expiresAt); err != nil {
		return err
	}

	slog.Info("password reset token created (mailer not implemented)",
		"email", user.Email,
		"token", raw,
		"expires_at", expiresAt,
	)
	return nil
}

func (s *AuthService) ConfirmPasswordReset(ctx context.Context, rawToken, newPassword string) error {
	hash := auth.HashResetToken(rawToken)
	token, err := s.resets.GetByHash(ctx, hash)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrInvalidResetToken
		}
		return err
	}
	if token.UsedAt != nil || time.Now().UTC().After(token.ExpiresAt) {
		return domain.ErrInvalidResetToken
	}

	passwordHash, err := auth.HashPassword(newPassword)
	if err != nil {
		return err
	}
	if err := s.users.UpdatePassword(ctx, token.UserID, passwordHash); err != nil {
		return err
	}
	if err := s.resets.MarkUsed(ctx, token.ID); err != nil {
		return err
	}
	_ = s.resets.InvalidateForUser(ctx, token.UserID)
	return nil
}

func (s *AuthService) issue(user domain.User) (AuthResult, error) {
	token, expiresAt, err := s.tokens.Issue(user.ID, user.Email, user.Role)
	if err != nil {
		return AuthResult{}, err
	}
	return AuthResult{User: user, Token: token, ExpiresAt: expiresAt}, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
