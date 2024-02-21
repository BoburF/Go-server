package routes

import (
	"crash-course-server/controllers"

	"github.com/gofiber/fiber/v2"
)

func AddUserGroup(app *fiber.App) {
	userGroup := app.Group("/user")

	userGroup.Post("/register", controllers.RegisterUser)
	userGroup.Post("/sign", controllers.SignIn)
}
