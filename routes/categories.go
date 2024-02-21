package routes

import (
	"crash-course-server/controllers"
	"crash-course-server/middleware"

	"github.com/gofiber/fiber/v2"
)

func AddCategoriesGroup(app *fiber.App) {
	categoriesGroup := app.Group("/categories")

	categoriesGroup.Use(middleware.AuthMiddleware)

	categoriesGroup.Post("/create", controllers.CreateCategory)
	categoriesGroup.Get("/", controllers.GetCategories)
}
