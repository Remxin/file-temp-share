package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	RoleUser  = "USER"
	RoleAdmin = "ADMIN"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Role         string
	IsActive     bool
	CreatedAt    time.Time
}
