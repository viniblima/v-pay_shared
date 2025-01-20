package healthcheck

import (
	"log"

	"github.com/streadway/amqp"
)

type HealthRabbitMQ interface {
	CheckRabbitMQ(rabbitmqURL string) error
}

type healthRabbitMQ struct{}

func (h *healthRabbitMQ) CheckRabbitMQ(rabbitmqURL string) error {

	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		log.Fatalf("error on connect rabbitmq: %v", err)
		return err
	}

	defer conn.Close()

	return nil
}

func NewHealthRabbitMQ() HealthRabbitMQ {
	return &healthRabbitMQ{}
}
