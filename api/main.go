package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ShaunakSensarma/URL-Shortener/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

// setupRoutes has the list of all routes.
func setupRoutes(app *fiber.App) {
	app.Get("/:url", routes.ResolveURL)
	app.Post("/api/v1", routes.ShortenURL)
}

// main function is the entry point of the code.
func main() {
	// /to load the environment variables defined in .env file.
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	app := fiber.New()

	app.Use(logger.New())

	setupRoutes(app)

	// to start the server.
	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}
