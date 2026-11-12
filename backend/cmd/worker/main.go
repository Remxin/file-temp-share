package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
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

	consumer, err := queue.NewConsumer(client)
	if err != nil {
		slog.Error("rabbitmq consumer", "err", err)
		os.Exit(1)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			slog.Error("consumer close", "err", err)
		}
	}()

	handler := func(ctx context.Context, msg queue.PasswordResetEmailMessage) error {
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

	slog.Info("worker started", "rabbitmq", cfg.RabbitMQURL, "from", cfg.GmailAppLogin)
	err = consumer.Run(ctx, handler)
	if err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("worker stopped", "err", err)
		os.Exit(1)
	}
	slog.Info("worker shut down")
}
