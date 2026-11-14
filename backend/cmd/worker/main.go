package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"file-temp-share/backend/internal/config"
	"file-temp-share/backend/internal/mailer"
	"file-temp-share/backend/internal/queue"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("config", "err", err)
		os.Exit(1)
	}
	if err := cfg.RequireGmail(); err != nil {
		slog.Error("config", "err", err)
		os.Exit(1)
	}

	gmail, err := mailer.NewGmailSMTP(cfg.GmailAppLogin, cfg.GmailAppPassword)
	if err != nil {
		slog.Error("mailer", "err", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	client, err := queue.Dial(cfg.RabbitMQURL)
	if err != nil {
		slog.Error("rabbitmq dial", "err", err)
		os.Exit(1)
	}
	defer func() {
		if err := client.Close(); err != nil {
			slog.Error("rabbitmq close", "err", err)
		}
	}()

	passwordResetConsumer, err := queue.NewConsumer(client)
	if err != nil {
		slog.Error("rabbitmq consumer", "err", err)
		os.Exit(1)
	}
	defer func() {
		if err := passwordResetConsumer.Close(); err != nil {
			slog.Error("password_reset consumer close", "err", err)
		}
	}()

	transferInviteConsumer, err := queue.NewConsumer(client)
	if err != nil {
		slog.Error("rabbitmq transfer_invite consumer", "err", err)
		os.Exit(1)
	}
	defer func() {
		if err := transferInviteConsumer.Close(); err != nil {
			slog.Error("transfer_invite consumer close", "err", err)
		}
	}()

	passwordResetHandler := func(ctx context.Context, msg queue.PasswordResetEmailMessage) error {
		slog.Info("sending password_reset email",
			"to", msg.To,
			"requested_at", msg.RequestedAt,
		)
		email, err := mailer.PasswordResetEmail(msg.To, msg.ResetURL)
		if err != nil {
			return err
		}
		return gmail.Send(ctx, email)
	}

	transferInviteHandler := func(ctx context.Context, msg queue.TransferInviteEmailMessage) error {
		slog.Info("sending transfer_invite email",
			"to", msg.To,
			"file", msg.FileName,
			"sent_at", msg.SentAt,
		)
		email, err := mailer.TransferInviteEmail(
			msg.To,
			msg.DownloadURL,
			msg.Title,
			msg.Message,
			msg.FileName,
			msg.FileSizeBytes,
			msg.FileCount,
			msg.ExpiresAt,
		)
		if err != nil {
			return err
		}
		return gmail.Send(ctx, email)
	}

	slog.Info("worker started", "rabbitmq", cfg.RabbitMQURL, "from", cfg.GmailAppLogin)

	var wg sync.WaitGroup
	errCh := make(chan error, 2)

	wg.Add(2)
	go func() {
		defer wg.Done()
		if runErr := passwordResetConsumer.Run(ctx, passwordResetHandler); runErr != nil {
			errCh <- runErr
		}
	}()
	go func() {
		defer wg.Done()
		if runErr := transferInviteConsumer.RunTransferInvite(ctx, transferInviteHandler); runErr != nil {
			errCh <- runErr
		}
	}()

	wg.Wait()
	close(errCh)

	for runErr := range errCh {
		if runErr != nil && !errors.Is(runErr, context.Canceled) {
			slog.Error("worker stopped", "err", runErr)
			os.Exit(1)
		}
	}
	slog.Info("worker shut down")
}
