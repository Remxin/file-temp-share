package mailer

import (
	"strings"
	"testing"
)

func TestPasswordResetEmailRendersFlashShareTemplate(t *testing.T) {
	msg, err := PasswordResetEmail("user@example.com", "https://flashshare.local/reset-password/confirm?token=abc")
	if err != nil {
		t.Fatalf("PasswordResetEmail: %v", err)
	}
	if msg.Subject != "Reset your FlashShare password" {
		t.Fatalf("subject: %q", msg.Subject)
	}
	if !strings.Contains(msg.HTML, "Reset Password") {
		t.Fatal("html missing CTA label")
	}
	if !strings.Contains(msg.HTML, "60 minutes") {
		t.Fatal("html missing expiry")
	}
	if !strings.Contains(msg.HTML, "https://flashshare.local/reset-password/confirm?token=abc") {
		t.Fatal("html missing reset url")
	}
	if !strings.Contains(msg.Text, "FlashShare") {
		t.Fatal("text missing brand")
	}
}
