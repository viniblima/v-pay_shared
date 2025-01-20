package messaging

import (
	"log"
	"sync"

	"github.com/streadway/amqp"
)

type RabbitMQCLient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

var (
	clientInstance *RabbitMQCLient
	once           sync.Once
)

// GetRabbitMQCLient returns a new instance of RabbitMQCLient
func GetRabbitMQCLient(url string) (*RabbitMQCLient, error) {
	once.Do(func() {
		conn, err := amqp.Dial(url)
		if err != nil {
			log.Fatalf("error on connect rabbitmq client: %v", err)
		}

		channel, err := conn.Channel()
		if err != nil {
			log.Fatalf("error on open a channel with rabbitmq client: %v", err)
		}

		clientInstance = &RabbitMQCLient{
			conn:    conn,
			channel: channel,
		}
	})
	return clientInstance, nil
}

// DeclareQueue creates a new queue
func (r *RabbitMQCLient) DeclareQueue(queueName string) (amqp.Queue, error) {
	return r.channel.QueueDeclare(queueName, true, false, false, false, nil)
}

// PublishMessage sends a message to a queue
func (r *RabbitMQCLient) PublishMessage(queueName string, body []byte) error {
	return r.channel.Publish("", queueName, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

// ConsumeMessage consumes a message from a queue
func (r *RabbitMQCLient) Close() {
	log.Println("Wallet Service is closing rabbitmq")
	r.channel.Close()
	r.conn.Close()
}
