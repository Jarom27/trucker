package messaging

import (
	"encoding/json"
	"fmt"
	"trucker/commands"

	"github.com/streadway/amqp"
)

type QueueMessenger struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

// Constructor para QueueMessenger
func NewQueueMessenger(queueName string, user string, pass string, host string, port string) (*QueueMessenger, error) {
	connection_string := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pass, host, port)
	conn, err := amqp.Dial(connection_string)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	_, err = channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	return &QueueMessenger{
		conn:    conn,
		channel: channel,
		queue:   queueName,
	}, nil
}

func (q *QueueMessenger) Send(data interface{}) error {
	switch value := data.(type) {
	case commands.LocationReport:
		location_json, err := json.Marshal(value)
		if err != nil {
			fmt.Printf("There was an error sending: %s", err)
		}
		q.publish(location_json)
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
		return fmt.Errorf("failed to publish a message: %w", err)
	}
	return nil
}
