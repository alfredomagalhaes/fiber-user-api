package routes

import (
	"github.com/alfredomagalhaes/fiber-user-api/controllers"
	"github.com/alfredomagalhaes/fiber-user-api/repositories"
	"github.com/gofiber/fiber/v2"
)

// UserRoutes defines all the routes that will deal with users context
func UserRoutes(route fiber.Router, repo repositories.UserRepository) {

	userGroup := route.Group("user")
	userGroup.Post("/", controllers.CreateNewUser(repo))
	userGroup.Get("/", controllers.GetUsers(repo))
}
