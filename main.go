package main

import (
	"crash-course-server/configs"
	"crash-course-server/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{"data": "Hello from Fiber and Mongodb"})
	})

	routes.AddUserGroup(app)
	routes.AddCategoriesGroup(app)
	routes.AddTransactionGroup(app)

	err := configs.LoadEnv()
	if err != nil {
		panic(err)
	}

	configs.ConnectDB()
	app.Listen(":4356")
}
