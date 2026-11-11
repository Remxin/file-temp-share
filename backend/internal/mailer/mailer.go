package mailer

import "context"

// Message is an outbound email. Prefer HTML + Text (multipart/alternative).
type Message struct {
	To      string
	Subject string
	Text    string // plain-text body (required)
	HTML    string // optional HTML body
}

// Mailer sends email messages.
type Mailer interface {
	Send(ctx context.Context, msg Message) error
}
