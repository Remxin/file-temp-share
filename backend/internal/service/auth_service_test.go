package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"file-temp-share/backend/internal/auth"
	"file-temp-share/backend/internal/domain"
	"file-temp-share/backend/internal/queue"
	"file-temp-share/backend/internal/repository"
)

type stubUsers struct {
	byID    map[uuid.UUID]domain.User
	byEmail map[string]domain.User
	updated map[uuid.UUID]string
}

func (s *stubUsers) Create(context.Context, string, string, string) (domain.User, error) {
	return domain.User{}, errors.New("not implemented")
}

func (s *stubUsers) GetByEmail(_ context.Context, email string) (domain.User, error) {
	u, ok := s.byEmail[email]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	return u, nil
}

func (s *stubUsers) GetByID(_ context.Context, id uuid.UUID) (domain.User, error) {
	u, ok := s.byID[id]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	return u, nil
}

func (s *stubUsers) UpdatePassword(_ context.Context, id uuid.UUID, passwordHash string) error {
	if s.updated == nil {
		s.updated = map[uuid.UUID]string{}
	}
	s.updated[id] = passwordHash
	if u, ok := s.byID[id]; ok {
		u.PasswordHash = passwordHash
		s.byID[id] = u
	}
	return nil
}

type stubResets struct {
	invalidated []uuid.UUID
}

func (s *stubResets) Create(context.Context, uuid.UUID, string, time.Time) (repository.PasswordResetToken, error) {
	return repository.PasswordResetToken{}, errors.New("not implemented")
}

func (s *stubResets) GetByHash(context.Context, string) (repository.PasswordResetToken, error) {
	return repository.PasswordResetToken{}, domain.ErrNotFound
}

func (s *stubResets) MarkUsed(context.Context, uuid.UUID) error { return nil }

func (s *stubResets) InvalidateForUser(_ context.Context, userID uuid.UUID) error {
	s.invalidated = append(s.invalidated, userID)
	return nil
}

type noopPublisher struct{}

func (noopPublisher) PublishPasswordReset(context.Context, queue.PasswordResetEmailMessage) error {
	return nil
}
func (noopPublisher) PublishTransferInvite(context.Context, queue.TransferInviteEmailMessage) error {
	return nil
}
func (noopPublisher) Close() error { return nil }

func newTestAuthService(users *stubUsers, resets *stubResets) *AuthService {
	tokens := auth.NewTokenManager("test-secret-for-unit-tests-32b!", time.Hour)
	return &AuthService{
		users:        users,
		resets:       resets,
		tokens:       tokens,
		publisher:    noopPublisher{},
		appPublicURL: "http://localhost:3000",
		resetTTL:     time.Hour,
	}
}

func TestLoginBusinessRules(t *testing.T) {
	hash, err := auth.HashPassword("correct-horse")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	activeID := uuid.New()
	inactiveID := uuid.New()
	users := &stubUsers{
		byID: map[uuid.UUID]domain.User{},
		byEmail: map[string]domain.User{
			"active@example.com": {
				ID:           activeID,
				Email:        "active@example.com",
				PasswordHash: hash,
				Role:         domain.RoleUser,
				IsActive:     true,
			},
			"blocked@example.com": {
				ID:           inactiveID,
				Email:        "blocked@example.com",
				PasswordHash: hash,
				Role:         domain.RoleUser,
				IsActive:     false,
			},
		},
	}
	svc := newTestAuthService(users, &stubResets{})

	t.Run("unknown email is invalid credentials", func(t *testing.T) {
		_, err := svc.Login(context.Background(), "missing@example.com", "x")
		if !errors.Is(err, domain.ErrInvalidCredentials) {
			t.Fatalf("got %v, want ErrInvalidCredentials", err)
		}
	})

	t.Run("inactive user", func(t *testing.T) {
		_, err := svc.Login(context.Background(), "blocked@example.com", "correct-horse")
		if !errors.Is(err, domain.ErrUserInactive) {
			t.Fatalf("got %v, want ErrUserInactive", err)
		}
	})

	t.Run("wrong password", func(t *testing.T) {
		_, err := svc.Login(context.Background(), "active@example.com", "nope")
		if !errors.Is(err, domain.ErrInvalidCredentials) {
			t.Fatalf("got %v, want ErrInvalidCredentials", err)
		}
	})
}

func TestChangePasswordBusinessRules(t *testing.T) {
	hash, err := auth.HashPassword("old-password")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	userID := uuid.New()
	users := &stubUsers{
		byID: map[uuid.UUID]domain.User{
			userID: {
				ID:           userID,
				Email:        "user@example.com",
				PasswordHash: hash,
				Role:         domain.RoleUser,
				IsActive:     true,
			},
		},
		byEmail: map[string]domain.User{},
	}
	resets := &stubResets{}
	svc := newTestAuthService(users, resets)

	t.Run("wrong current password", func(t *testing.T) {
		err := svc.ChangePassword(context.Background(), userID, "bad", "new-password")
		if !errors.Is(err, domain.ErrValidation) {
			t.Fatalf("got %v, want ErrValidation", err)
		}
	})

	t.Run("new password must differ", func(t *testing.T) {
		err := svc.ChangePassword(context.Background(), userID, "old-password", "old-password")
		if !errors.Is(err, domain.ErrValidation) {
			t.Fatalf("got %v, want ErrValidation", err)
		}
	})

	t.Run("inactive user", func(t *testing.T) {
		inactiveID := uuid.New()
		users.byID[inactiveID] = domain.User{
			ID:           inactiveID,
			PasswordHash: hash,
			IsActive:     false,
		}
		err := svc.ChangePassword(context.Background(), inactiveID, "old-password", "new-password")
		if !errors.Is(err, domain.ErrUserInactive) {
			t.Fatalf("got %v, want ErrUserInactive", err)
		}
	})

	t.Run("success updates password and invalidates resets", func(t *testing.T) {
		err := svc.ChangePassword(context.Background(), userID, "old-password", "brand-new")
		if err != nil {
			t.Fatalf("ChangePassword: %v", err)
		}
		if _, ok := users.updated[userID]; !ok {
			t.Fatal("expected password update")
		}
		if len(resets.invalidated) == 0 || resets.invalidated[len(resets.invalidated)-1] != userID {
			t.Fatalf("expected reset invalidation for %s, got %v", userID, resets.invalidated)
		}
	})
}

func TestRequestPasswordResetNoEnumeration(t *testing.T) {
	svc := newTestAuthService(&stubUsers{
		byID:    map[uuid.UUID]domain.User{},
		byEmail: map[string]domain.User{},
	}, &stubResets{})

	if err := svc.RequestPasswordReset(context.Background(), "nobody@example.com"); err != nil {
		t.Fatalf("missing user should succeed silently: %v", err)
	}
}
