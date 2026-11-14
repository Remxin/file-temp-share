package mailer

import (
	"fmt"
	"math"
	"time"
)

type transferInviteData struct {
	Title         string
	Message       string
	FileName      string
	FileSizeLabel string
	FileKindLabel string
	DownloadURL   string
	ExpiresInLabel string
}

// TransferInviteEmail renders the FlashShare transfer-invitation templates.
func TransferInviteEmail(
	to string,
	downloadURL string,
	title, message string,
	fileName string,
	fileSizeBytes int64,
	fileCount int,
	expiresAt time.Time,
) (Message, error) {
	data := transferInviteData{
		Title:          title,
		Message:        message,
		FileName:       fileName,
		FileSizeLabel:  formatByteSize(fileSizeBytes),
		FileKindLabel:  fileKindLabel(fileCount),
		DownloadURL:    downloadURL,
		ExpiresInLabel: formatExpiresIn(time.Until(expiresAt)),
	}

	text, err := renderText("transfer_invite.txt", data)
	if err != nil {
		return Message{}, err
	}
	html, err := renderHTML("transfer_invite.html", data)
	if err != nil {
		return Message{}, err
	}

	subject := "A file is ready for download on FlashShare"
	if title != "" {
		subject = fmt.Sprintf("%s — ready to download on FlashShare", title)
	}

	return Message{
		To:      to,
		Subject: subject,
		Text:    text,
		HTML:    html,
	}, nil
}

func fileKindLabel(fileCount int) string {
	if fileCount <= 1 {
		return "Single Document"
	}
	return fmt.Sprintf("%d Documents", fileCount)
}

func formatByteSize(n int64) string {
	if n < 1024 {
		return fmt.Sprintf("%d B", n)
	}
	units := []string{"KB", "MB", "GB", "TB"}
	v := float64(n)
	unit := "B"
	for _, u := range units {
		v /= 1024
		unit = u
		if v < 1024 {
			break
		}
	}
	if v >= 100 || math.Mod(v, 1) < 0.05 {
		return fmt.Sprintf("%.0f %s", v, unit)
	}
	return fmt.Sprintf("%.1f %s", v, unit)
}

func formatExpiresIn(d time.Duration) string {
	if d <= 0 {
		return "0 hours"
	}
	hours := int(math.Round(d.Hours()))
	if hours < 1 {
		mins := int(math.Round(d.Minutes()))
		if mins <= 1 {
			return "1 minute"
		}
		return fmt.Sprintf("%d minutes", mins)
	}
	if hours == 1 {
		return "1 hour"
	}
	if hours%24 == 0 {
		days := hours / 24
		if days == 1 {
			return "1 day"
		}
		return fmt.Sprintf("%d days", days)
	}
	return fmt.Sprintf("%d hours", hours)
}
