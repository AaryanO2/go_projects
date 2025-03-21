package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

var mg MongoInstance

const dbName = "fiber-hrms"

type Employee struct {
	ID     string  `json:"id,omitempty" bson:"_id,omitempty" `
	Name   string  `json:"name"`
	Salary float64 `json:"salary"`
	Age    float64 `json:"age"`
}

func Connect() error {
	mongoUser := os.Getenv("mongo_user")
	mongoPassword := os.Getenv("mongo_password")
	mongoHost := "mongodb"
	mongoPort := "27017"

	mongoURL := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		mongoUser, mongoPassword, mongoHost, mongoPort)
	log.Printf("Connecting to MongoDB at: %s", mongoHost)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	defer cancel()
	if err != nil {
		return err
	}
	db := client.Database(dbName)
	mg = MongoInstance{
		Client: client,
		Db:     db,
	}
	return nil
}

func main() {
	app := fiber.New()
	if err := Connect(); err != nil {
		fmt.Println("oo guru")
		log.Fatal(err)
	}
	app.Get("/employee", func(c fiber.Ctx) error {
		query := bson.D{{}}
		cursor, err := mg.Db.Collection("employees").Find(c.Context(), query)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		var employees []Employee = make([]Employee, 0)
		if err := cursor.All(c.Context(), &employees); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		return c.JSON(employees)
	})
	app.Post("/employee", func(c fiber.Ctx) error {
		collection := mg.Db.Collection("employees")
		employee := new(Employee)
		if err := c.Bind().JSON(employee); err != nil {
			return c.Status(400).SendString("Error parsing Body: " + err.Error())
		}
		employee.ID = ""

		insertionResult, err := collection.InsertOne(c.Context(), employee)

		if err != nil {
			return c.Status(500).SendString("Error while inserting: " + err.Error())
		}

		filter := bson.D{{Key: "_id", Value: insertionResult.InsertedID}}

		createdRecord := collection.FindOne(c.Context(), filter)

		createdEmployee := &Employee{}

		if err = createdRecord.Decode(createdEmployee); err != nil {
			return c.Status(500).SendString("Error decoding created record: " + err.Error())
		}

		return c.Status(201).JSON(createdEmployee)

	})
	app.Put("/employee/:id", func(c fiber.Ctx) error {
		idParam := c.Params("id")
		employeeId, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			return c.SendStatus(400)
		}
		employee := new(Employee)

		if err = c.Bind().JSON(employee); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		query := bson.D{{Key: "_id", Value: employeeId}}
		update := bson.D{
			{
				Key: "$set",
				Value: bson.D{
					{Key: "name", Value: employee.Name},
					{Key: "age", Value: employee.Age},
					{Key: "salary", Value: employee.Salary},
				},
			},
		}

		err = mg.Db.Collection("employees").FindOneAndUpdate(c.Context(), query, update).Err()
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.SendStatus(404)
			}
			return c.SendStatus(500)
		}

		employee.ID = idParam

		return c.Status(200).JSON(employee)

	})

	app.Delete("/employee/:id", func(c fiber.Ctx) error {
		id := c.Params("id")
		employeeId, err := primitive.ObjectIDFromHex(id)

		if err != nil {
			return c.SendStatus(400)
		}

		query := bson.D{{
			Key:   "_id",
			Value: employeeId,
		}}

		result, err := mg.Db.Collection("employees").DeleteOne(c.Context(), query)

		if err != nil {
			return c.SendStatus(500)
		}
		if result.DeletedCount < 1 {
			return c.SendStatus(404)
		}

		return c.Status(200).JSON(fiber.Map{"Message": "record deleted"})

	})
	log.Fatal(app.Listen(":3000"))
}
