package messaging

import (
	"encoding/json"
	"fmt"
	"time"
	"trucker/commands"

	"github.com/streadway/amqp"
)

type QueueMessenger struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
	user    string
	pass    string
	host    string
	port    string
}

func NewQueueMessenger(queueName, user, pass, host, port string) (*QueueMessenger, error) {
	connection_string := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pass, host, port)
	var conn *amqp.Connection
	var channel *amqp.Channel
	var err error

	const maxRetries = 5
	for i := 1; i <= maxRetries; i++ {
		conn, err = amqp.Dial(connection_string)
		if err == nil {
			break
		}
		fmt.Printf("Attempt %d: Failed to connect to RabbitMQ: %v\n", i, err)
		time.Sleep(time.Duration(1<<i) * time.Second) // Espera exponencial (1s, 2s, 4s, 8s, 16s)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ after %d attempts: %w", maxRetries, err)
	}

	channel, err = conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	// Declarar la cola
	_, err = channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	fmt.Println("Successfully connected to RabbitMQ")

	return &QueueMessenger{
		conn:    conn,
		channel: channel,
		queue:   queueName,
		user:    user,
		pass:    pass,
		host:    host,
		port:    port,
	}, nil
}

func (q *QueueMessenger) Send(data interface{}) error {
	switch value := data.(type) {
	case commands.LocationReport:
		location_json, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal data: %w", err)
		}
		err = q.publish(location_json)
		if err != nil {
			return fmt.Errorf("failed to send message to queue: %w", err)
		}
	default:
		return fmt.Errorf("unsupported data type: %T", value)
	}
	return nil
}

func (q *QueueMessenger) publish(data []byte) error {
	err := q.channel.Publish(
		"",
		q.queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)
	if err != nil {
		fmt.Println("Publish failed, attempting to reconnect...")

		reconnectErr := q.reconnect()
		if reconnectErr != nil {
			return fmt.Errorf("failed to reconnect to RabbitMQ: %w", reconnectErr)
		}

		// Intentar publicar nuevamente después de reconectar
		return q.channel.Publish(
			"",
			q.queue,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        data,
			},
		)
	}
	return nil
}

func (q *QueueMessenger) reconnect() error {
	fmt.Println("Reconnecting to RabbitMQ...")

	var conn *amqp.Connection
	var channel *amqp.Channel
	var err error

	const maxRetries = 5
	for i := 1; i <= maxRetries; i++ {
		connection_string := fmt.Sprintf("amqp://%s:%s@%s:%s/", q.user, q.pass, q.host, q.port)
		conn, err = amqp.Dial(connection_string)
		if err == nil {
			break
		}
		fmt.Printf("Reconnect attempt %d failed: %v\n", i, err)
		time.Sleep(time.Duration(1<<i) * time.Second)
	}

	if err != nil {
		return fmt.Errorf("failed to reconnect to RabbitMQ after %d attempts: %w", maxRetries, err)
	}

	channel, err = conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel after reconnection: %w", err)
	}

	q.conn = conn
	q.channel = channel

	fmt.Println("Reconnected to RabbitMQ successfully")
	return nil
}

func (q *QueueMessenger) Close() {
	if q.channel != nil {
		q.channel.Close()
	}
	if q.conn != nil {
		q.conn.Close()
	}
	fmt.Println("RabbitMQ connection closed")
}
