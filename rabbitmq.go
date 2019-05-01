package database

import (
	"encoding/json"
	"time"

	"github.com/streadway/amqp"
)

// Message represents message received from listener for processing
type Message struct {
	MessageID string
	Timestamp time.Time
	Body      []byte
}

// ListenerAction passes back message body to your provided callback for processing
type ListenerAction func(Message, func(bool) error, func(bool, bool) error)

// RabbitMQConnect gets rabbitmq connection
func RabbitMQConnect(address string) (*amqp.Connection, error) {
	return amqp.Dial(address)
}

// ListenToQueue creates a non blocking listener to rabbitmq
func ListenToQueue(conn *amqp.Connection, queue string, action ListenerAction) error {
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		queue,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			action(Message{MessageID: d.MessageId, Timestamp: d.Timestamp, Body: d.Body}, d.Ack, d.Nack)
		}
	}()

	<-forever
	// return never hits till we close channel forever
	return nil
}

// SendToQueue broadcasts payload to queue
func SendToQueue(conn *amqp.Connection, queue string, payload interface{}) error {
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()
	q, err := ch.QueueDeclare(
		queue,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	p, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        p,
		})
}
