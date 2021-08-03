package routes

import (
	"fmt"
	"exp-clean-arch-arangodb/api/middlewares"
	"exp-clean-arch-arangodb/pkg/activity"
	"exp-clean-arch-arangodb/pkg/city"
	"exp-clean-arch-arangodb/pkg/entities"
	"exp-clean-arch-arangodb/pkg/session"
	"exp-clean-arch-arangodb/pkg/user"
	"exp-clean-arch-arangodb/utils"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

const PasswordTooShort = "Le mot de passe doit être égal ou supérieur à 9 caractères"
const WrongCredendials = "Ces informations ne correspondent à aucun utilisateur"
const ConfirmMsg = "Un email vient de vous être envoyer afin de confirmer votre compte"
const CityNotFound = "Cette commune n'existe pas"

func UserRouter(app fiber.Router, userService user.Service, sessionService session.Service, cityService city.Service, activityService activity.Service) {
	app.Post("/auth", authUser(userService, sessionService))
	app.Post("/enreg", createUser(userService))

	mustBeAuth := app.Group("a", middlewares.Authorization(sessionService), middlewares.SetLocals(), middlewares.UserLastActivity(activityService))

	mustBeAuth.Post("/ajo_c", addCityUser(userService, sessionService, cityService))
	mustBeAuth.Post("/del_c", deleteCityUser(userService, sessionService))
	mustBeAuth.Post("/del_u", deleteUser(userService, sessionService))
}

func authUser(userService user.Service, sessionService session.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var authUserBody entities.AuthUserBody

		if err := c.BodyParser(&authUserBody); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"msg": err,
			})
		}

		email := authUserBody.Email
		password := authUserBody.Password

		email = strings.TrimSpace(email)
		password = strings.TrimSpace(password)

		err := utils.EmailValidity(email)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"msg": err.Error(),
			})
		}

		user, err := userService.FindByEmail(email)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"msg": err.Error(),
			})
		}

		match, err := argon2id.ComparePasswordAndHash(password, user.Hpwd)
		if !match {
			return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
				"msg": WrongCredendials,
			})
		} else if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"msg": err.Error(),
			})
		}

		sessionID, err := sessionService.Set(user.Key)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"msg": err.Error(),
			})
		}

		cookie := fasthttp.Cookie{}
		cookie.SetKey("sid")
		cookie.SetValue(fmt.Sprintf("%s;HttpOnly", sessionID))
		c.Context().Response.Header.SetCookie(&cookie)
		c.Context().Response.Header.SetBytesV("Access-Control-Allow-Origin", c.Context().Request.Header.Peek("Origin"))

		return c.SendStatus(fiber.StatusAccepted)
	}
}

func createUser(userService user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var newUserBody entities.NewUserBody

		if err := c.BodyParser(&newUserBody); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"msg": err,
			})
		}

		username := newUserBody.Username
		email := newUserBody.Email
		password := newUserBody.Password

		username = strings.TrimSpace(username)
		email = strings.TrimSpace(email)
		password = strings.TrimSpace(password)

		// vérifie si l'username est valide
		if err := utils.UsernameValidity(username); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(&fiber.Map{
				"msg": err,
			})
		}

		// vérifie si l'adresse email est au bon format
		if err := utils.EmailValidity(email); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(&fiber.Map{
				"msg": err,
			})
		}

		if len(newUserBody.Password) < 9 {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(&fiber.Map{
				"msg": PasswordTooShort,
			})
		}
		newUserBody.Username = username
		newUserBody.Email = email
		newUserBody.Password = password

		username, err := userService.Create(newUserBody)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"msg": err,
			})
		}

		return c.JSON(&fiber.Map{
			"msg": fmt.Sprintf("L'utilisateur %s a bien été créé.\n\nVous receverez très prochainement un email afin de confirmer votre compte.\nSi celui-ci n'est pas confirmé dans les 2 semaines, le compte sera automatiquement bloqué.", username),
		})
	}
}

func addCityUser(userService user.Service, sessionService session.Service, cityService city.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var citySearchBody entities.CitySearchBody

		if err := c.BodyParser(&citySearchBody); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"msg": err,
			})
		}

		cityFound, err := cityService.IsCityExist(citySearchBody.Nom)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"msg": CityNotFound,
			})
		}

		userID := c.Locals("userID").(string)
		err = userService.AddCity(userID, cityFound.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"msg": err.Error(),
			})
		}

		return c.SendStatus(fiber.StatusAccepted)
	}
}

func deleteCityUser(userService user.Service, sessionService session.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(string)

		err := userService.DeleteCity(userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"msg": err.Error(),
			})
		}

		return c.SendStatus(fiber.StatusAccepted)
	}
}

func deleteUser(userService user.Service, sessionService session.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionID := c.Locals("sessionID").(string)
		userID := c.Locals("userID").(string)

		if err := userService.DeleteTotallyAccount(userID); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"msg": err.Error(),
			})
		}

		if err := sessionService.Delete(sessionID); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"msg": err.Error(),
			})
		}

		return c.JSON(&fiber.Map{
			"msg": "Votre compte ainsi que les informations ont bien été supprimées",
		})
	}
}
