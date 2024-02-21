package middleware

import (
	"crash-course-server/configs"
	"crash-course-server/models"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	token := c.Cookies("Authorization")

	if token == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Unauthorized",
			"message": "Token not found",
		})
	}

	sessionColl := configs.GetCollection("session")
	var session models.Session
	err := sessionColl.FindOne(c.Context(), fiber.Map{"token": token}).Decode(&session)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Unauthorized",
			"message": "Token not found",
		})
	}

	c.Locals("userId", session.UserID)

	return c.Next()
}
