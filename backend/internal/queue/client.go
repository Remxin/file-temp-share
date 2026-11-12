package queue

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Client owns a single AMQP connection. Channels are created per Publisher/Consumer.
type Client struct {
	conn *amqp.Connection
}

// Dial opens a connection to RabbitMQ.
func Dial(url string) (*Client, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("amqp dial: %w", err)
	}
	return &Client{conn: conn}, nil
}

// Conn returns the underlying connection (for advanced use / tests).
func (c *Client) Conn() *amqp.Connection {
	return c.conn
}

// Close closes the connection. Safe to call once; subsequent use of derived
// publishers/consumers will fail.
func (c *Client) Close() error {
	if c == nil || c.conn == nil || c.conn.IsClosed() {
		return nil
	}
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("amqp close: %w", err)
	}
	return nil
}

// openChannel creates a channel and declares topology on it.
func (c *Client) openChannel() (*amqp.Channel, error) {
	if c == nil || c.conn == nil {
		return nil, fmt.Errorf("amqp client is nil")
	}
	ch, err := c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("open channel: %w", err)
	}
	if err := Declare(ch); err != nil {
		_ = ch.Close()
		return nil, err
	}
	return ch, nil
}
