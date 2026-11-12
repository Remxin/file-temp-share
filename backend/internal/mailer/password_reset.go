package mailer

// PasswordResetExpireMinutes matches AuthService password-reset token TTL (1h).
const PasswordResetExpireMinutes = 60

type passwordResetData struct {
	ResetURL      string
	ExpireMinutes int
}

// PasswordResetEmail renders the FlashShare password-reset templates.
func PasswordResetEmail(to, resetURL string) (Message, error) {
	data := passwordResetData{
		ResetURL:      resetURL,
		ExpireMinutes: PasswordResetExpireMinutes,
	}

	text, err := renderText("password_reset.txt", data)
	if err != nil {
		return Message{}, err
	}
	html, err := renderHTML("password_reset.html", data)
	if err != nil {
		return Message{}, err
	}

	return Message{
		To:      to,
		Subject: "Reset your FlashShare password",
		Text:    text,
		HTML:    html,
	}, nil
}
