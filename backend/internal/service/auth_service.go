package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"

	"file-temp-share/backend/internal/auth"
	"file-temp-share/backend/internal/domain"
	"file-temp-share/backend/internal/queue"
	"file-temp-share/backend/internal/repository"
)

const passwordResetTTL = time.Hour

// userStore is the AuthService subset of user persistence (keeps unit tests free of sqlc).
type userStore interface {
	Create(ctx context.Context, email, passwordHash, role string) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error
}

type passwordResetStore interface {
	Create(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) (repository.PasswordResetToken, error)
	GetByHash(ctx context.Context, tokenHash string) (repository.PasswordResetToken, error)
	MarkUsed(ctx context.Context, id uuid.UUID) error
	InvalidateForUser(ctx context.Context, userID uuid.UUID) error
}

type AuthService struct {
	users        userStore
	resets       passwordResetStore
	tokens       *auth.TokenManager
	publisher    queue.Publisher
	appPublicURL string
	resetTTL     time.Duration
}

func NewAuthService(
	users userStore,
	resets passwordResetStore,
	tokens *auth.TokenManager,
	publisher queue.Publisher,
	appPublicURL string,
) *AuthService {
	return &AuthService{
		users:        users,
		resets:       resets,
		tokens:       tokens,
		publisher:    publisher,
		appPublicURL: strings.TrimRight(appPublicURL, "/"),
		resetTTL:     passwordResetTTL,
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

// RequestPasswordReset always succeeds from the caller's perspective when the
// user is missing/inactive (no email enumeration). When a token is created,
// a password-reset email job is enqueued.
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

	resetURL := fmt.Sprintf("%s/reset-password/confirm?token=%s", s.appPublicURL, raw)
	msg := queue.PasswordResetEmailMessage{
		Type:        queue.TypePasswordReset,
		To:          user.Email,
		ResetURL:    resetURL,
		RequestedAt: time.Now().UTC(),
	}
	if err := s.publisher.PublishPasswordReset(ctx, msg); err != nil {
		return fmt.Errorf("enqueue password reset email: %w", err)
	}

	slog.Info("password reset email enqueued",
		"email", user.Email,
		"expires_at", expiresAt,
	)
	return nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if !user.IsActive {
		return domain.ErrUserInactive
	}
	if !auth.CheckPassword(user.PasswordHash, currentPassword) {
		// 400 (not 401): wrong current password must not clear the session JWT.
		return fmt.Errorf("%w: current password is incorrect", domain.ErrValidation)
	}
	if currentPassword == newPassword {
		return fmt.Errorf("%w: new password must differ from current password", domain.ErrValidation)
	}

	passwordHash, err := auth.HashPassword(newPassword)
	if err != nil {
		return err
	}
	if err := s.users.UpdatePassword(ctx, user.ID, passwordHash); err != nil {
		return err
	}
	_ = s.resets.InvalidateForUser(ctx, user.ID)
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
