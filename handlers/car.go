package handlers

import (
	"github.com/cmerin0/SimpleCarsApp/db"
	"github.com/cmerin0/SimpleCarsApp/models"
	"github.com/gofiber/fiber/v2"
)

func GetCars(c *fiber.Ctx) error {
	cars := []models.Car{} // Creates an array of Models cars

	db.DB.Db.Find(&cars) // Find all the cars and store them in the cars variable

	return c.Render("cars", fiber.Map{
		"Title": "Cars",
		"Cars":  cars,
	})
}

func GetCarById(c *fiber.Ctx) error {
	id := c.Params("id") // Get the id from the request

	car := new(models.Car) // Create a new car

	db.DB.Db.First(&car, id) // Find the car by id and store it in the car variable

	if car.ID == 0 { // Check if the car exists
		return c.Status(404).JSON("Car not found") // Return a status 404 code response and the error
	}

	return c.Status(200).JSON(car) // Return a status 200 code response and the car
}

func CreateCar(c *fiber.Ctx) error {
	car := new(models.Car) // Create a new car

	// Check if body sent is properly sent
	if err := c.BodyParser(car); err != nil { // Parse the body of the request and store it in the car variable
		return c.Status(400).JSON(err.Error()) // Return a status 400 code response and the error
	}

	db.DB.Db.Create(&car) // Create the car in the database

	return c.Status(201).JSON(car) // Return a status 201 code response and the car
}

func UpdateCar(c *fiber.Ctx) error {
	id := c.Params("id")   // Get the id from the request
	car := new(models.Car) // Create a new car object

	// Check if body sent is properly sent
	if err := c.BodyParser(car); err != nil { // Parse the body of the request and store it in the car variable
		return c.Status(400).JSON(err.Error()) // Return a status 400 code response and the error
	}

	// Check if the Car exists
	existingCar := new(models.Car)
	result := db.DB.Db.First(&existingCar, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Car not found",
		})
	}

	// Update the Car
	car.ID = existingCar.ID // Ensure the ID is set for the update
	db.DB.Db.Save(&car)

	return c.Status(200).JSON(car) // return a status 200 response + Car updated
}

func DeleteCar(c *fiber.Ctx) error {
	id := c.Params("id") // Get the id from the request

	car := new(models.Car) // Create a new car object

	result := db.DB.Db.First(&car, id) // Find the car by id and store it in the car variable

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Car not found",
		})
	}

	db.DB.Db.Delete(&car) // Delete the car from the database

	return c.SendStatus(204) // return a status 204 response
}
