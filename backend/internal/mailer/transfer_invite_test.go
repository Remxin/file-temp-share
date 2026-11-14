package mailer

import (
	"strings"
	"testing"
	"time"
)

func TestTransferInviteEmailRendersFlashShareTemplate(t *testing.T) {
	msg, err := TransferInviteEmail(
		"recv@example.com",
		"https://flashshare.local/t/abc123",
		"Lecture notes",
		"Please review before Friday.",
		"lecture-notes-final.pdf",
		24_600_000,
		1,
		time.Now().UTC().Add(48*time.Hour),
	)
	if err != nil {
		t.Fatalf("TransferInviteEmail: %v", err)
	}
	if !strings.Contains(msg.Subject, "Lecture notes") {
		t.Fatalf("subject: %q", msg.Subject)
	}
	if !strings.Contains(msg.HTML, "Transfer name") {
		t.Fatal("html missing transfer name label")
	}
	if !strings.Contains(msg.HTML, "Lecture notes") {
		t.Fatal("html missing title")
	}
	if !strings.Contains(msg.HTML, "Please review before Friday.") {
		t.Fatal("html missing message")
	}
	if !strings.Contains(msg.HTML, "lecture-notes-final.pdf") {
		t.Fatal("html missing filename")
	}
	if !strings.Contains(msg.HTML, "Go to Download Page") {
		t.Fatal("html missing CTA")
	}
	if !strings.Contains(msg.HTML, "https://flashshare.local/t/abc123") {
		t.Fatal("html missing download url")
	}
	if !strings.Contains(msg.Text, "Transfer name: Lecture notes") {
		t.Fatal("text missing title")
	}
}

func TestTransferInviteEmailOmitsEmptyTitleAndMessage(t *testing.T) {
	msg, err := TransferInviteEmail(
		"recv@example.com",
		"https://flashshare.local/t/abc123",
		"",
		"",
		"notes.pdf",
		1024,
		1,
		time.Now().UTC().Add(time.Hour),
	)
	if err != nil {
		t.Fatalf("TransferInviteEmail: %v", err)
	}
	if strings.Contains(msg.HTML, "Transfer name") {
		t.Fatal("html should omit empty title section")
	}
	if strings.Contains(msg.HTML, ">Message<") {
		t.Fatal("html should omit empty message section")
	}
	if msg.Subject != "A file is ready for download on FlashShare" {
		t.Fatalf("subject: %q", msg.Subject)
	}
}

func TestFormatByteSize(t *testing.T) {
	if got := formatByteSize(24_600_000); got != "23.5 MB" {
		t.Fatalf("formatByteSize: %q", got)
	}
}
