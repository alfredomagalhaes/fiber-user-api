package routes

import (
	"github.com/alfredomagalhaes/fiber-user-api/controllers"
	"github.com/gofiber/fiber/v2"
)

func LoginRoutes(router fiber.Router) {

	//Create a new group to login endpoints
	login := router.Group("login")
	login.Post("/", controllers.LoginUser)
}
