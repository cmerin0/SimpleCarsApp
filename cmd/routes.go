package main

import (
	"github.com/cmerin0/SimpleCarsApp/db"
	"github.com/cmerin0/SimpleCarsApp/handlers"
	"github.com/gofiber/fiber/v2"
)

// Function that gathers the routes in a single file
func setupRoutes(app *fiber.App) {

	// Main route of test
	app.Get("/", handlers.Home)

	// Routes of Authentication
	auth := app.Group("/auth", handlers.VerifyToken)
	auth.Get("/users", handlers.GetUsers)
	auth.Post("/logout", handlers.Logout)
	app.Post("/login", handlers.Login)
	app.Post("/register", handlers.Register)

	// Routes of Make
	makesRoute := app.Group("/makes")
	makesRoute.Get("/", func(c *fiber.Ctx) error {
		return handlers.GetMakes(c, db.Cache.RedisClient)
	})
	makesRoute.Get("/:id", handlers.GetMakeById)
	makesRoute.Post("/", handlers.CreateMake)
	makesRoute.Put("/:id", handlers.UpdateMake)
	makesRoute.Delete("/:id", handlers.DeleteMake)

	// Routes of Car
	carsRoute := app.Group("/cars")
	carsRoute.Get("/", handlers.GetCars)
	carsRoute.Get("/:id", handlers.GetCarById)
	carsRoute.Post("/", handlers.CreateCar)
	carsRoute.Put("/:id", handlers.UpdateCar)
	carsRoute.Delete("/:id", handlers.DeleteCar)

}
