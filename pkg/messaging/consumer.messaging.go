package messaging

import (
	"context"

	"github.com/streadway/amqp"
)

type RabbitMQConsumer struct {
	client *RabbitMQCLient
}

func NewRabbitMQConsumer(client *RabbitMQCLient) *RabbitMQConsumer {
	return &RabbitMQConsumer{
		client: client,
	}
}

func (r *RabbitMQConsumer) Start(queueName string, contextLocal context.Context, handler func(msg amqp.Delivery, contextLocal context.Context)) error {
	msgs, err := r.client.channel.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			handler(msg, contextLocal)
		}
	}()

	return nil
}
