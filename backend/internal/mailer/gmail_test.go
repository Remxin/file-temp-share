package mailer

import "testing"

func TestFormatFromHeader(t *testing.T) {
	got := formatFromHeader("FlashShare", "noreply@example.com")
	want := "FlashShare <noreply@example.com>"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}

	if got := formatFromHeader("", "a@b.c"); got != "a@b.c" {
		t.Fatalf("empty name: got %q", got)
	}
}
