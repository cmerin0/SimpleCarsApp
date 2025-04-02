package main

import (
	"github.com/cmerin0/SimpleCarsApp/db"
	"github.com/gofiber/fiber/v2"
)

func main() {

	// Connect to the database
	db.ConnectDB()

	// Connect to the cache database
	db.ConnectCache()

	// Initializing app with Fiber new
	app := fiber.New()

	// Importing routes from routes package
	setupRoutes(app)

	// App listening in port 3000
	app.Listen(":3000")
}
