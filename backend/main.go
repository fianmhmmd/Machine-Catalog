package main

import (
	"log"
	"os"

	"github.com/fianmhmmd/machine-catalog/backend/database"
	"github.com/fianmhmmd/machine-catalog/backend/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect to database
	database.ConnectDB()

	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("ALLOWED_ORIGINS"),
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Routes
	routes.SetupRoutes(app)

	// Health check
	app.Get("/api/v1/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status": "success",
			"message": "API is healthy",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(app.Listen(":" + port))
}
