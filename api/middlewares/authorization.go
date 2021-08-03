package middlewares

import (
	"exp-clean-arch-arangodb/pkg/session"

	"github.com/gofiber/fiber/v2"
)

func Authorization(sessionService session.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		sidCookie := c.Cookies("sid")

		if sidCookie != "" {
			if err := sessionService.Get(sidCookie); err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
					"msg": err.Error(),
				})
			}

			return c.Next()
		}

		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"msg": "Cette session n'existe pas. Veuillez vous authentifier",
		})
	}
}
