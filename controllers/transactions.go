package controllers

import (
	"crash-course-server/configs"
	"crash-course-server/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type transactionDto struct {
	CategoryId string `json:"categoryId" bson:"categoryId"`
	Amount     int    `json:"amount" bson:"amount"`
}

func CreateTransaction(c *fiber.Ctx) error {
	b := new(transactionDto)
	if err := c.BodyParser(b); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid body",
			"message": err,
		})
	}

	transactionsColl := configs.GetCollection("transactions")

	userID, ok := c.Locals("userId").(string)
	if !ok {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": "Failed to retrieve userID from context",
		})
	}

	transaction := models.Transactions{
		CreatedTime: time.Now().In(time.FixedZone("UTC+5", 5*60*60)),
		CategoryId:  b.CategoryId,
		UserId:      userID,
		Amount:      b.Amount,
	}

	insertResult, err := transactionsColl.InsertOne(c.Context(), transaction)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create transaction",
			"message": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{"result": fiber.Map{"transactionId": insertResult.InsertedID}})
}

func GetTransactions(c *fiber.Ctx) error {
	transactionsColl := configs.GetCollection("transactions")

	cursor, err := transactionsColl.Find(c.Context(), bson.D{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to get transactions",
			"message": err.Error(),
		})
	}

	var transactions []models.Transactions
	if err = cursor.All(c.Context(), &transactions); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to find transactions",
			"message": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{"result": fiber.Map{"transactions": transactions}})
}

func GetTransactionsByWeek(c *fiber.Ctx) error {
	weekNumber, err := strconv.Atoi(c.Params("weekNumber"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid week number",
			"message": err.Error(),
		})
	}

	now := time.Now()
	lastWeekStart := now.AddDate(0, 0, -int(now.Weekday())-(weekNumber-1)*7).Truncate(24 * time.Hour)
	lastWeekEnd := now.AddDate(0, 0, -int(now.Weekday())-(weekNumber-1)*7+6).Truncate(24 * time.Hour).Add(time.Hour*23 + time.Minute*59 + time.Second*59)

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"createdTime": bson.M{
					"$gte": lastWeekStart,
					"$lte": lastWeekEnd,
				},
			},
		},
		{
			"$lookup": bson.M{
				"from": "categories",
				"let":  bson.M{"categoryId": "$categoryId"},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"$expr": bson.M{
								"$eq": []interface{}{"$_id", bson.M{"$toObjectId": "$$categoryId"}},
							},
						},
					},
				},
				"as": "category",
			},
		},
		{
			"$unwind": "$category",
		},
		{
			"$group": bson.M{
				"_id":          nil,
				"transactions": bson.M{"$push": "$$ROOT"},
				"totalMinusAmount": bson.M{
					"$sum": bson.M{
						"$cond": bson.M{
							"if":   bson.M{"$lt": []interface{}{"$amount", 0}},
							"then": "$amount",
							"else": 0,
						},
					},
				},
				"totalPlusAmount": bson.M{
					"$sum": bson.M{
						"$cond": bson.M{
							"if":   bson.M{"$gt": []interface{}{"$amount", 0}},
							"then": "$amount",
							"else": 0,
						},
					},
				},
				"overallAmount": bson.M{"$sum": "$amount"},
			},
		},
		{
			"$project": bson.M{
				"_id":              0,
				"transactions":     1,
				"totalMinusAmount": 1,
				"totalPlusAmount":  1,
				"overallAmount":    1,
			},
		},
	}

	transactionColl := configs.GetCollection("transactions")

	cursor, err := transactionColl.Aggregate(c.Context(), pipeline)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to find transactions",
			"message": err.Error(),
		})
	}

	var transactions []bson.M
	if err := cursor.All(c.Context(), &transactions); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to find transactions",
			"message": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{"result": fiber.Map{"transactions": transactions}})
}
