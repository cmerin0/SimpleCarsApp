package main

import (
	"github.com/cmerin0/SimpleCarsApp/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {

	// Connect to the database
	db.ConnectDB()

	// Initializing Fiber template engine
	engine := html.New("./views", ".html")

	// Initializing app with Fiber new
	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
	})

	// Importing routes from routes package
	setupRoutes(app)

	// App listening in port 3000
	app.Listen(":3000")
}
