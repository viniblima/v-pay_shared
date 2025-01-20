package messaging

import "github.com/streadway/amqp"

type RabbitMQConsumer struct {
	client *RabbitMQCLient
}

func NewRabbitMQConsumer(client *RabbitMQCLient) *RabbitMQConsumer {
	return &RabbitMQConsumer{
		client: client,
	}
}

func (r *RabbitMQConsumer) Start(queueName string, handler func(msg amqp.Delivery)) error {
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
		print(msgs)
		for msg := range msgs {
			handler(msg)
		}
	}()

	return nil
}
