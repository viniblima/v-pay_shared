package queries

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"v-pay_shared/healthcheck"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HealthQuery interface {
	Execute(ctx *fiber.Ctx) error
}

type healthQuery struct {
	dbHealth       healthcheck.HealthDatabase
	rabbitMqHealth healthcheck.HealthRabbitMQ
}

func (q *healthQuery) Execute(ctx *fiber.Ctx) error {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if err := q.rabbitMqHealth.CheckRabbitMQ(rabbitURL); err != nil {
		log.Printf("error on connect rabbitmq of wallet service: %v", err)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "unhealthy",
			"message": "RabbitMQ is unavailable",
		})
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Shanghai",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	if err := q.dbHealth.CheckDatabase(dsn); err != nil {
		log.Printf("error on connect database of wallet service: %v", err)
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "unhealthy",
			"message": "Database is unavailable",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"status":  "healthy",
		"message": "Services is running smoothly",
	})
}

func NewHealthQuery(db *gorm.DB) HealthQuery {
	return &healthQuery{
		dbHealth:       healthcheck.NewHealthDatabase(db),
		rabbitMqHealth: healthcheck.NewHealthRabbitMQ(),
	}
}
