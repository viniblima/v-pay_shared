package middlewares

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/viniblima/v-pay_shared/pkg/enums"
	"github.com/viniblima/v-pay_shared/pkg/models"
)

type cacheMiddleware struct {
	redisClient  redis.Client
	contextLocal context.Context
}

type CacheMiddleware interface {
	VerifyWallet(ctx *fiber.Ctx) error
}

func NewCacheMiddleware(redisClient redis.Client, contextLocal context.Context) CacheMiddleware {
	return &cacheMiddleware{
		redisClient:  redisClient,
		contextLocal: contextLocal,
	}
}

func (m *cacheMiddleware) VerifyWallet(ctx *fiber.Ctx) error {

	userID := ctx.Locals("userID").(string)
	list, err := m.redisClient.Get(m.contextLocal, enums.CacheListUserWallet).Result()

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "error on get wallet list",
		})
	}

	obj := []byte(list)
	listWallet := []models.Wallet{}

	err = json.Unmarshal(obj, &listWallet)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "error on parse list",
		})
	}

	filtered := []models.Wallet{}
	for _, item := range listWallet {
		if item.User == userID {
			filtered = append(filtered, item)
		}
	}

	if len(filtered) > 0 {
		return ctx.Next()
	} else {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "User is not associated to this wallet",
		})
	}
}
