package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

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
	VerifyUserWallet(ctx *fiber.Ctx) error
	VerifyStoreWallet(ctx *fiber.Ctx) error
}

func NewCacheMiddleware(redisClient redis.Client, contextLocal context.Context) CacheMiddleware {
	return &cacheMiddleware{
		redisClient:  redisClient,
		contextLocal: contextLocal,
	}
}

func (m *cacheMiddleware) VerifyStoreWallet(ctx *fiber.Ctx) error {
	sID := ctx.Get("Store")

	if sID != "" {
		split := strings.Split(sID, "key=")

		if len(split) < 2 {
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid tag",
			})
		}

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

		filtered := models.Wallet{}
		found := false
		for _, item := range listWallet {
			if item.User == userID && item.ID == split[1] {
				filtered = item
				found = true
				break
			}
		}

		if found {
			ctx.Locals("walletID", filtered.ID)
			return ctx.Next()
		} else {
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "User do not have a wallet",
			})
		}
	}
	return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
		"error": "Invalid headers",
	})
}

func (m *cacheMiddleware) VerifyUserWallet(ctx *fiber.Ctx) error {

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

	filtered := models.Wallet{}
	found := false
	for _, item := range listWallet {
		if item.User == userID {
			filtered = item
			found = true
			break
		}
	}

	if found {
		ctx.Locals("walletID", filtered.ID)
		return ctx.Next()
	} else {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "User is not associated to this wallet",
		})
	}
}
