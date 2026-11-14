package queue

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	ExchangeEmails = "filetemp.emails"

	QueuePasswordReset      = "email.password_reset"
	RoutingKeyPasswordReset = "email.password_reset"

	QueueTransferInvite      = "email.transfer_invite"
	RoutingKeyTransferInvite = "email.transfer_invite"

	prefetchCount = 1
)

// Declare ensures the emails topic exchange, email queues, and bindings exist.
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

	if err := declareBoundQueue(ch, QueuePasswordReset, RoutingKeyPasswordReset); err != nil {
		return err
	}
	if err := declareBoundQueue(ch, QueueTransferInvite, RoutingKeyTransferInvite); err != nil {
		return err
	}

	return nil
}

func declareBoundQueue(ch *amqp.Channel, queueName, routingKey string) error {
	if _, err := ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	); err != nil {
		return fmt.Errorf("declare queue %q: %w", queueName, err)
	}

	if err := ch.QueueBind(
		queueName,
		routingKey,
		ExchangeEmails,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("bind queue %q to %q: %w", queueName, ExchangeEmails, err)
	}
	return nil
}
