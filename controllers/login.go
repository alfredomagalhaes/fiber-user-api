package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func LoginUser(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(fiber.Map{"success": true, "message": "user logged in", "token": "123456"})
}
