package routes

import (
	"crash-course-server/controllers"
	"crash-course-server/middleware"

	"github.com/gofiber/fiber/v2"
)

func AddTransactionGroup(app *fiber.App) {
	transactionGroup := app.Group("/transactions")

	transactionGroup.Use(middleware.AuthMiddleware)

	transactionGroup.Post("/create", controllers.CreateTransaction)
	transactionGroup.Get("/", controllers.GetTransactions)
	transactionGroup.Get("/:weekNumber", controllers.GetTransactionsByWeek)
}
