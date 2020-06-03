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
type ListenerAction func(*Message) error

// RabbitMQConnect gets rabbitmq connection
func RabbitMQConnect(address string) (*amqp.Connection, error) {
	return amqp.Dial(address)
}

// ListenToQueue creates a non blocking listener to rabbitmq
func ListenToQueue(conn *amqp.Connection, queue string, durable bool, listenChan chan bool, errorChan chan error, action ListenerAction) {
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		errorChan <- err
		return
	}

	q, err := ch.QueueDeclare(
		queue,
		durable,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		errorChan <- err
		return
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		errorChan <- err
		return
	}

	go func() {
		for d := range msgs {
			processingErr := action(&Message{MessageID: d.MessageId, Timestamp: d.Timestamp, Body: d.Body})
			if processingErr != nil {
				errorChan <- processingErr
			}
		}
	}()

	// block
	<-listenChan
}

// SendToQueue broadcasts payload to queue
func SendToQueue(conn *amqp.Connection, queue string, durable bool, payload interface{}) error {
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()
	q, err := ch.QueueDeclare(
		queue,
		durable,
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
