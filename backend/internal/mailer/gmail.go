package mailer

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"
)

const (
	gmailSMTPHost = "smtp.gmail.com"
	gmailSMTPAddr = "smtp.gmail.com:587"
)

// GmailSMTP sends mail via Gmail SMTP using an app password.
type GmailSMTP struct {
	username string // full Gmail address (also From)
	password string // Google app password
}

// NewGmailSMTP builds a Mailer for Gmail SMTP (STARTTLS on :587).
func NewGmailSMTP(login, appPassword string) (*GmailSMTP, error) {
	login = strings.TrimSpace(login)
	appPassword = strings.TrimSpace(appPassword)
	if login == "" {
		return nil, fmt.Errorf("GMAIL_APP_LOGIN is required")
	}
	if appPassword == "" {
		return nil, fmt.Errorf("GMAIL_APP_PASSWORD is required")
	}
	return &GmailSMTP{username: login, password: appPassword}, nil
}

// Send delivers the message through Gmail SMTP (multipart/alternative when HTML is set).
func (g *GmailSMTP) Send(ctx context.Context, msg Message) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	to := strings.TrimSpace(msg.To)
	if to == "" {
		return fmt.Errorf("mailer: recipient is empty")
	}
	if strings.TrimSpace(msg.Subject) == "" {
		return fmt.Errorf("mailer: subject is empty")
	}
	if strings.TrimSpace(msg.Text) == "" {
		return fmt.Errorf("mailer: text body is empty")
	}

	payload, err := buildMIME(g.username, to, msg)
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", g.username, g.password, gmailSMTPHost)

	errCh := make(chan error, 1)
	go func() {
		errCh <- smtp.SendMail(gmailSMTPAddr, auth, g.username, []string{to}, payload)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("gmail smtp send: %w", err)
		}
		return nil
	}
}

func buildMIME(from, to string, msg Message) ([]byte, error) {
	var b strings.Builder
	b.WriteString("From: " + from + "\r\n")
	b.WriteString("To: " + to + "\r\n")
	b.WriteString("Subject: " + sanitizeHeader(msg.Subject) + "\r\n")
	b.WriteString("MIME-Version: 1.0\r\n")

	if strings.TrimSpace(msg.HTML) == "" {
		b.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		b.WriteString("\r\n")
		b.WriteString(msg.Text)
		return []byte(b.String()), nil
	}

	const boundary = "flashshare-alt-boundary"
	b.WriteString("Content-Type: multipart/alternative; boundary=" + boundary + "\r\n")
	b.WriteString("\r\n")

	b.WriteString("--" + boundary + "\r\n")
	b.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	b.WriteString("Content-Transfer-Encoding: 8bit\r\n")
	b.WriteString("\r\n")
	b.WriteString(msg.Text)
	b.WriteString("\r\n")

	b.WriteString("--" + boundary + "\r\n")
	b.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	b.WriteString("Content-Transfer-Encoding: 8bit\r\n")
	b.WriteString("\r\n")
	b.WriteString(msg.HTML)
	b.WriteString("\r\n")

	b.WriteString("--" + boundary + "--\r\n")
	return []byte(b.String()), nil
}

func sanitizeHeader(s string) string {
	return strings.Map(func(r rune) rune {
		if r == '\r' || r == '\n' {
			return -1
		}
		return r
	}, s)
}
