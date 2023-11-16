package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/alfredomagalhaes/fiber-user-api/repositories"
	"github.com/alfredomagalhaes/fiber-user-api/types"
	"github.com/gofiber/fiber/v2"
)

func CreateNewUser(repo repositories.UserRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var userRequest types.User

		err := c.BodyParser(&userRequest)

		if err != nil {
			log.Printf("%v", err)
			return c.Status(http.StatusBadRequest).JSON(ErrorResponse(errors.New("could not parse the request")))
		}

		err = repo.Save(&userRequest)

		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(ErrorResponse(errors.New("could not create a new user")))
		}

		return c.Status(http.StatusCreated).JSON(fiber.Map{
			"message": "user created",
		})
	}

}

func GetUsers(repo repositories.UserRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var users []types.User

		users, err := repo.ListAll()

		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(ErrorResponse(errors.New("could not list users")))
		}

		return c.Status(http.StatusCreated).JSON(fiber.Map{
			"message": "",
			"data":    users,
		})
	}
}
