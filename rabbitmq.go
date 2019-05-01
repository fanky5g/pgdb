package database

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

// ListenerAction passes back message body to your provided callback for processing
type ListenerAction func(messageID string, timestamp string, body []byte)

// RabbitMQConnect gets rabbitmq connection
func RabbitMQConnect(address string) *amqp.Connection {
	conn, err := amqp.Dial(address)
	HandleError(err)
	return conn
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
			action(d.MessageId, d.Timestamp, d.Body)
		}
	}()

	<-forever
}

// SendToQueue broadcasts payload to queue
func SendToQueue(conn *amqp.Connection, queue string, payload interface{}) error {
	defer redisConn.Close()

	ch, err := redisConn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()
	q, err := ch.QueueDeclare(
		name,
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
