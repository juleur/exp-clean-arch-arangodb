package routes

import (
	"exp-clean-arch-arangodb/api/middlewares"
	"exp-clean-arch-arangodb/pkg/city"
	"exp-clean-arch-arangodb/pkg/entities"
	"exp-clean-arch-arangodb/pkg/session"

	"github.com/gofiber/fiber/v2"
)

func CityRouter(app fiber.Router, sessionService session.Service, cityService city.Service) {
	mustBeAuth := app.Group("a", middlewares.Authorization(sessionService))

	mustBeAuth.Post("trouver_nom_commune", searchCity(cityService))
}

func searchCity(cityService city.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var citySearchBody entities.CitySearchBody

		if err := c.BodyParser(&citySearchBody); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"msg": err,
			})
		}

		citiesName, err := cityService.FindCity(citySearchBody.Nom)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"msg": err.Error(),
			})
		}

		return c.JSON(citiesName)
	}
}
