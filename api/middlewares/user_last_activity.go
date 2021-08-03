package middlewares

import (
	"exp-clean-arch-arangodb/pkg/activity"

	"github.com/gofiber/fiber/v2"
)

func UserLastActivity(activityService activity.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		sidCookie := c.Cookies("sid")

		activityService.LastUserActivity(sidCookie)

		return c.Next()
	}
}
