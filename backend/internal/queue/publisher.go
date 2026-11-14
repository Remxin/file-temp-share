package queue

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Publisher publishes email jobs to RabbitMQ.
type Publisher interface {
	PublishPasswordReset(ctx context.Context, msg PasswordResetEmailMessage) error
	PublishTransferInvite(ctx context.Context, msg TransferInviteEmailMessage) error
	Close() error
}

// RabbitPublisher publishes to the emails topic exchange.
type RabbitPublisher struct {
	ch *amqp.Channel
}

// NewPublisher opens a dedicated channel, declares topology, and returns a Publisher.
func NewPublisher(client *Client) (*RabbitPublisher, error) {
	ch, err := client.openChannel()
	if err != nil {
		return nil, fmt.Errorf("publisher: %w", err)
	}
	return &RabbitPublisher{ch: ch}, nil
}

// PublishPasswordReset publishes a durable password-reset email job.
func (p *RabbitPublisher) PublishPasswordReset(ctx context.Context, msg PasswordResetEmailMessage) error {
	if msg.Type == "" {
		msg.Type = TypePasswordReset
	}
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal password_reset message: %w", err)
	}

	err = p.ch.PublishWithContext(
		ctx,
		ExchangeEmails,
		RoutingKeyPasswordReset,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("publish password_reset: %w", err)
	}
	return nil
}

// PublishTransferInvite publishes a durable transfer-invitation email job.
func (p *RabbitPublisher) PublishTransferInvite(ctx context.Context, msg TransferInviteEmailMessage) error {
	if msg.Type == "" {
		msg.Type = TypeTransferInvite
	}
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal transfer_invite message: %w", err)
	}

	err = p.ch.PublishWithContext(
		ctx,
		ExchangeEmails,
		RoutingKeyTransferInvite,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("publish transfer_invite: %w", err)
	}
	return nil
}

// Close closes the publisher channel.
func (p *RabbitPublisher) Close() error {
	if p == nil || p.ch == nil {
		return nil
	}
	if err := p.ch.Close(); err != nil {
		return fmt.Errorf("publisher close: %w", err)
	}
	return nil
}
