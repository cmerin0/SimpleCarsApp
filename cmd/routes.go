package main

import (
	"github.com/cmerin0/SimpleCarsApp/handlers"
	"github.com/gofiber/fiber/v2"
)

// Function that gathers the routes in a single file
func setupRoutes(app *fiber.App) {

	// Main route of test
	app.Get("/", handlers.Home)

	// Routes of Authentication
	auth := app.Group("/auth", handlers.VerifyToken)
	app.Get("/users", handlers.GetUsers)
	auth.Post("/logout", handlers.Logout)
	app.Post("/login", handlers.Login)
	app.Post("/register", handlers.Register)

	// Routes of Make
	app.Get("/makes", handlers.GetMakes)
	app.Get("/makes/:id", handlers.GetMakeById)
	app.Post("/makes", handlers.CreateMake)
	app.Put("/makes/:id", handlers.UpdateMake)
	app.Delete("/makes/:id", handlers.DeleteMake)

	// Routes of Car
	app.Get("/cars", handlers.GetCars)
	app.Get("/cars/:id", handlers.GetCarById)
	app.Post("/cars", handlers.CreateCar)
	app.Put("/cars/:id", handlers.UpdateCar)
	app.Delete("/cars/:id", handlers.DeleteCar)

}
