package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cmerin0/SimpleCarsApp/db"
	"github.com/cmerin0/SimpleCarsApp/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/redis/go-redis/v9"
)

const cacheKeyAllMakes = "all_makes"
const cacheExpiration = 1 * time.Minute

var ctx = context.Background()

func Home(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"message": "Welcome to the Simple Cars App",
	})
}

// Route that returns all the makes
func GetMakes(c *fiber.Ctx, redisClient *redis.Client) error {

	// Check if data exists in Redis cache
	cachedMakes, err := redisClient.Get(ctx, cacheKeyAllMakes).Result()

	if err == nil {
		log.Info("Serving from Redis Cache")
		var makes []models.Make
		if err := json.Unmarshal([]byte(cachedMakes), &makes); err != nil {
			log.Fatal("Error unmarshalling cached makes: %v\n", err) // Fallback to database if cache data is corrupted
			return fetchMakesFromDBAndCache(c, redisClient)
		}
		// Returning cached makes data
		return c.Status(fiber.StatusOK).JSON(makes)
	}

	if err != redis.Nil {
		fmt.Printf("Error accessing Redis cache: %v\n", err) // Don't immediately fail, fallback to database
	}

	return fetchMakesFromDBAndCache(c, redisClient) // Data not in cache or Redis error, fetch from database
}

func fetchMakesFromDBAndCache(c *fiber.Ctx, redisClient *redis.Client) error {

	makes := []models.Make{}                              // Creates an array of Models make
	db.DB.Db.Order("ID asc").Preload("Cars").Find(&makes) // Find all the makes and store them in the makes variable

	// Cache the fetched data in Redis
	makesJSON, err := json.Marshal(makes)

	if err != nil {
		fmt.Printf("Error marshalling makes for cache: %v\n", err)
		return c.Status(fiber.StatusOK).JSON(makes) // Continue without caching if marshalling fails
	}

	// Setting the makes (Inserting into Redis)
	err = redisClient.Set(ctx, cacheKeyAllMakes, makesJSON, cacheExpiration).Err()

	if err != nil {
		fmt.Printf("Error setting makes in Redis cache: %v\n", err)
	}

	log.Info("Serving from Database")
	return c.Status(fiber.StatusOK).JSON(makes)
}

// Route that returns a make by ID
func GetMakeById(c *fiber.Ctx) error {

	id := c.Params("id")     // Get the id from the url
	make := new(models.Make) // Create a new model variable

	db.DB.Db.Preload("Cars").Find(&make, id) // Find the make by the id and store it in the make variable

	// Throwing error if ID is not found
	if make.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Make not found",
		})
	}

	return c.Status(200).JSON(make) // return a status 200 response and the make by ID
}

// Route that creates a new Make field
func CreateMake(c *fiber.Ctx) error {

	make := new(models.Make) // Create a new model variable

	// Body parser returns an error && if Not error
	if err := c.BodyParser(make); err != nil {
		// We return an internal server error then a map that in Json Format
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// If not errors we create the new make and send a status 200 code response and the make created
	db.DB.Db.Create(&make)          //package.global_variable.instace.pointer_to_gorm_database_connection
	return c.Status(201).JSON(make) // Returns status 201 Created + Make created

}

func UpdateMake(c *fiber.Ctx) error {

	fmt.Println("entro aca al menos")

	id := c.Params("id")     // Get the id from the url
	make := new(models.Make) // Create a new model variable

	// Body parser returns an error if body is not valid
	if err := c.BodyParser(make); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Check if the make exists
	existingMake := new(models.Make)
	result := db.DB.Db.First(&existingMake, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Make not found",
		})
	}

	// Update the make
	make.ID = existingMake.ID // Ensure the ID is set for the update
	db.DB.Db.Save(&make)

	return c.Status(200).JSON(make) // return a status 200 Ok response + make Updated

}

func DeleteMake(c *fiber.Ctx) error {

	id := c.Params("id") // Get the id from the url

	// Check if the make exists
	existingMake := new(models.Make)
	result := db.DB.Db.First(&existingMake, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Make not found",
		})
	}

	// Delete the make
	db.DB.Db.Delete(&existingMake)

	return c.SendStatus(204) // return a status 204 response No Content

}
