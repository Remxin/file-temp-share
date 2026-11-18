package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"file-temp-share/backend/internal/domain"
	"file-temp-share/backend/internal/repository"
)

type AdminService struct {
	users        *repository.UserRepository
	transfers    *repository.TransferRepository
	appPublicURL string
}

func NewAdminService(
	users *repository.UserRepository,
	transfers *repository.TransferRepository,
	appPublicURL string,
) *AdminService {
	return &AdminService{
		users:        users,
		transfers:    transfers,
		appPublicURL: strings.TrimRight(appPublicURL, "/"),
	}
}

func (s *AdminService) ListUsers(ctx context.Context) ([]domain.User, error) {
	return s.users.List(ctx)
}

func (s *AdminService) SetUserActive(ctx context.Context, actorID, targetID uuid.UUID, isActive bool) (domain.User, error) {
	if actorID == targetID && !isActive {
		return domain.User{}, fmt.Errorf("%w: cannot deactivate yourself", domain.ErrForbidden)
	}

	target, err := s.users.GetByID(ctx, targetID)
	if err != nil {
		return domain.User{}, err
	}

	if target.IsActive == isActive {
		return target, nil
	}

	if !isActive && target.Role == domain.RoleAdmin {
		n, err := s.users.CountActiveAdmins(ctx)
		if err != nil {
			return domain.User{}, err
		}
		if n <= 1 {
			return domain.User{}, fmt.Errorf("%w: cannot deactivate the last admin", domain.ErrForbidden)
		}
	}

	return s.users.SetActive(ctx, targetID, isActive)
}

func (s *AdminService) ListTransfers(ctx context.Context) ([]domain.AdminTransferListItem, error) {
	items, err := s.transfers.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	for i := range items {
		items[i].Status = items[i].EffectiveStatus(now)
	}
	return items, nil
}

func (s *AdminService) SetTransferStatus(ctx context.Context, transferID uuid.UUID, status string) (domain.Transfer, error) {
	switch status {
	case domain.TransferStatusActive, domain.TransferStatusBlocked:
	default:
		return domain.Transfer{}, fmt.Errorf("%w: status must be active or blocked", domain.ErrValidation)
	}

	transfer, err := s.transfers.GetByID(ctx, transferID)
	if err != nil {
		return domain.Transfer{}, err
	}
	if transfer.IsDeleted() {
		return domain.Transfer{}, domain.ErrNotFound
	}
	if transfer.Status == status {
		return transfer, nil
	}

	updated, err := s.transfers.UpdateStatus(ctx, transferID, status)
	if err != nil {
		return domain.Transfer{}, err
	}
	updated.Status = updated.EffectiveStatus(time.Now().UTC())
	return updated, nil
}

func (s *AdminService) ShareURL(t domain.Transfer) string {
	return s.appPublicURL + "/t/" + t.ShortCode
}
