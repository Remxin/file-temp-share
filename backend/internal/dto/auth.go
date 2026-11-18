package dto

import (
	"time"

	"github.com/google/uuid"

	"file-temp-share/backend/internal/domain"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type PasswordResetConfirmRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=72"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8,max=72"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type AuthResponse struct {
	Token     string       `json:"token"`
	ExpiresAt time.Time    `json:"expires_at"`
	User      UserResponse `json:"user"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

func NewUserResponse(u domain.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
	}
}
