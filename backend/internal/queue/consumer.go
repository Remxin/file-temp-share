package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

// PasswordResetHandler processes a single password-reset email job.
type PasswordResetHandler func(ctx context.Context, msg PasswordResetEmailMessage) error

// Consumer consumes password-reset jobs from RabbitMQ.
type Consumer struct {
	ch *amqp.Channel
}

// NewConsumer opens a dedicated channel, declares topology, and configures QoS.
func NewConsumer(client *Client) (*Consumer, error) {
	ch, err := client.openChannel()
	if err != nil {
		return nil, fmt.Errorf("consumer: %w", err)
	}
	if err := ch.Qos(prefetchCount, 0, false); err != nil {
		_ = ch.Close()
		return nil, fmt.Errorf("consumer qos: %w", err)
	}
	return &Consumer{ch: ch}, nil
}

// Run consumes password-reset messages until ctx is cancelled or the channel closes.
// Invalid JSON is nack'd without requeue. Handler errors are nack'd without requeue.
func (c *Consumer) Run(ctx context.Context, handler PasswordResetHandler) error {
	if handler == nil {
		return fmt.Errorf("password reset handler is nil")
	}

	deliveries, err := c.ch.Consume(
		QueuePasswordReset,
		"",    // consumer tag — server-generated
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		return fmt.Errorf("consume %q: %w", QueuePasswordReset, err)
	}

	slog.Info("consuming password_reset email jobs",
		"queue", QueuePasswordReset,
		"exchange", ExchangeEmails,
		"routing_key", RoutingKeyPasswordReset,
	)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case d, ok := <-deliveries:
			if !ok {
				return fmt.Errorf("amqp deliveries channel closed")
			}
			c.handleDelivery(ctx, d, handler)
		}
	}
}

func (c *Consumer) handleDelivery(ctx context.Context, d amqp.Delivery, handler PasswordResetHandler) {
	var msg PasswordResetEmailMessage
	if err := json.Unmarshal(d.Body, &msg); err != nil {
		slog.Error("invalid password_reset message json",
			"err", err,
			"body", string(d.Body),
		)
		if nackErr := d.Nack(false, false); nackErr != nil {
			slog.Error("nack poison message", "err", nackErr)
		}
		return
	}

	if err := handler(ctx, msg); err != nil {
		slog.Error("password_reset handler failed",
			"err", err,
			"to", msg.To,
		)
		if nackErr := d.Nack(false, false); nackErr != nil {
			slog.Error("nack failed job", "err", nackErr)
		}
		return
	}

	if err := d.Ack(false); err != nil {
		slog.Error("ack password_reset job", "err", err)
	}
}

// Close closes the consumer channel (cancels active consumers).
func (c *Consumer) Close() error {
	if c == nil || c.ch == nil {
		return nil
	}
	if err := c.ch.Close(); err != nil {
		return fmt.Errorf("consumer close: %w", err)
	}
	return nil
}
