package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func VerifyJWT(c *fiber.Ctx) error {
	auth := c.Get("Authorization")
	claims := jwt.MapClaims{}
	if auth != "" {
		split := strings.Split(auth, "JWT ")

		if len(split) < 2 {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid tag",
			})
		}

		jwt.ParseWithClaims(split[1], claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("PASSWORD_SECRET")), nil
		})

		if claims["sub"] != "auth" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		c.Locals("userID", claims["ID"])

		return c.Next()
	} else {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid headers",
		})
	}
}
