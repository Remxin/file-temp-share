package domain

import "errors"

var (
	ErrNotFound            = errors.New("not found")
	ErrConflict            = errors.New("conflict")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrUserInactive        = errors.New("user inactive")
	ErrValidation          = errors.New("validation failed")
	ErrInvalidResetToken   = errors.New("invalid or expired reset token")
	ErrTransferExpired     = errors.New("transfer expired")
	ErrTransferBlocked     = errors.New("transfer blocked")
)
