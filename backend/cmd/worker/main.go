package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"file-temp-share/backend/internal/config"
	"file-temp-share/backend/internal/queue"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("config", "err", err)
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

	slog.Info("worker started", "rabbitmq", cfg.RabbitMQURL)
	err = consumer.Run(ctx, logPasswordReset)
	if err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("worker stopped", "err", err)
		os.Exit(1)
	}
	slog.Info("worker shut down")
}

func logPasswordReset(_ context.Context, msg queue.PasswordResetEmailMessage) error {
	slog.Info("received password_reset email job",
		"type", msg.Type,
		"to", msg.To,
		"reset_url", msg.ResetURL,
		"requested_at", msg.RequestedAt,
	)
	return nil
}
