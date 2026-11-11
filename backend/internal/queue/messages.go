package queue

import "time"

const TypePasswordReset = "password_reset"

// PasswordResetEmailMessage is the JSON body for password-reset email jobs.
type PasswordResetEmailMessage struct {
	Type        string    `json:"type"`
	To          string    `json:"to"`
	ResetURL    string    `json:"reset_url"`
	RequestedAt time.Time `json:"requested_at"`
}
