package middlewares

import (
	"exp-clean-arch-arangodb/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func SetLocals() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		sidCookie := c.Cookies("sid")

		split := strings.Split(sidCookie, "u")
		userID, err := utils.DecodeString(split[0])
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"msg": "Une erreur est survenue avec votre session",
			})
		}

		c.Locals("userID", userID)
		c.Locals("sessionID", c.Cookies("sid"))

		return c.Next()
	}
}
