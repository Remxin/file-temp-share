package queue

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	ExchangeEmails = "filetemp.emails"

	QueuePasswordReset      = "email.password_reset"
	RoutingKeyPasswordReset = "email.password_reset"

	prefetchCount = 1
)

// Declare ensures the emails topic exchange, password-reset queue, and binding exist.
// Safe to call on every process start (idempotent).
func Declare(ch *amqp.Channel) error {
	if err := ch.ExchangeDeclare(
		ExchangeEmails,
		"topic",
		true,  // durable
		false, // auto-delete
		false, // internal
		false, // no-wait
		nil,
	); err != nil {
		return fmt.Errorf("declare exchange %q: %w", ExchangeEmails, err)
	}

	if _, err := ch.QueueDeclare(
		QueuePasswordReset,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	); err != nil {
		return fmt.Errorf("declare queue %q: %w", QueuePasswordReset, err)
	}

	if err := ch.QueueBind(
		QueuePasswordReset,
		RoutingKeyPasswordReset,
		ExchangeEmails,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("bind queue %q to %q: %w", QueuePasswordReset, ExchangeEmails, err)
	}

	return nil
}
