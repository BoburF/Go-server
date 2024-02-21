package controllers

import (
	"crash-course-server/configs"
	"crash-course-server/models"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type userDto struct {
	Name    string `json:"name" bson:"name"`
	Surname string `json:"surname" bson:"surname"`
}

func RegisterUser(c *fiber.Ctx) error {
	b := new(userDto)
	if err := c.BodyParser(b); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid body",
			"message": err,
		})
	}

	userColl := configs.GetCollection("user")

	result := userColl.FindOne(c.Context(), bson.M{"name": b.Name, "surname": b.Surname})
	if result.Err() == nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "User already exists",
			"message": "A user with the same name and surname already exists.",
		})
	}

	insertResult, err := userColl.InsertOne(c.Context(), b)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create user",
			"message": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{"result": fiber.Map{"userId": insertResult.InsertedID}})
}

func SignIn(c *fiber.Ctx) error {
	b := new(userDto)
	if err := c.BodyParser(b); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid body",
			"message": err,
		})
	}

	userColl := configs.GetCollection("user")

	result := userColl.FindOne(c.Context(), bson.M{"name": b.Name, "surname": b.Surname})
	if result.Err() != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to found User",
			"message": result.Err(),
		})
	}

	var user models.User
	result.Decode(&user)

	token := generateUniqueToken()

	newSession := models.Session{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().In(time.FixedZone("UTC+5", 5*60*60)),
	}

	sessionColl := configs.GetCollection("session")

	_, err := sessionColl.InsertOne(c.Context(), newSession)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create session",
			"message": err.Error(),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:   "Authorization",
		Value:  token,
		MaxAge: 900,
	})

	return c.Status(200).JSON(fiber.Map{
		"message": "Sign-in successful",
		"token":   token,
	})
}

func generateUniqueToken() string {
	randomBytes := make([]byte, 32)
	rand.Read(randomBytes)

	token := base64.URLEncoding.EncodeToString(randomBytes)

	return token
}
