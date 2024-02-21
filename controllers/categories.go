package controllers

import (
	"crash-course-server/configs"
	"crash-course-server/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type categoryDto struct {
	Title string `json:"title" bson:"title"`
}

func CreateCategory(c *fiber.Ctx) error {
	b := new(categoryDto)
	if err := c.BodyParser(b); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid body",
			"message": err,
		})
	}

	categoryColl := configs.GetCollection("categories")

	result := categoryColl.FindOne(c.Context(), bson.M{"title": b.Title})
	if result.Err() == nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Category already exists",
			"message": result.Err(),
		})
	}

	insertResult, err := categoryColl.InsertOne(c.Context(), b)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create category",
			"message": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{"result": fiber.Map{"categoryId": insertResult.InsertedID}})
}

func GetCategories(c *fiber.Ctx) error {
	categoryColl := configs.GetCollection("categories")

	cursor, err := categoryColl.Find(c.Context(), bson.D{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to get categories",
			"message": err.Error(),
		})
	}

	var categories []models.Categories
	if err = cursor.All(c.Context(), &categories); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to find categories",
			"message": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{"result": fiber.Map{"categories": categories}})
}
