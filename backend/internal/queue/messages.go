package queue

import "time"

const (
	TypePasswordReset   = "password_reset"
	TypeTransferInvite  = "transfer_invite"
)

// PasswordResetEmailMessage is the JSON body for password-reset email jobs.
type PasswordResetEmailMessage struct {
	Type        string    `json:"type"`
	To          string    `json:"to"`
	ResetURL    string    `json:"reset_url"`
	RequestedAt time.Time `json:"requested_at"`
}

// TransferInviteEmailMessage is the JSON body for transfer-invitation email jobs.
type TransferInviteEmailMessage struct {
	Type          string    `json:"type"`
	To            string    `json:"to"`
	DownloadURL   string    `json:"download_url"`
	Title         string    `json:"title,omitempty"`
	Message       string    `json:"message,omitempty"`
	FileName      string    `json:"file_name"`
	FileSizeBytes int64     `json:"file_size_bytes"`
	FileCount     int       `json:"file_count"`
	ExpiresAt     time.Time `json:"expires_at"`
	SentAt        time.Time `json:"sent_at"`
}
