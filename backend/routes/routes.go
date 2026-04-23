package routes

import (
	"github.com/fianmhmmd/machine-catalog/backend/handlers"
	"github.com/fianmhmmd/machine-catalog/backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/login", handlers.Login)
	auth.Post("/refresh", handlers.Refresh)

	// Public routes
	// api.Get("/categories", handlers.GetCategories)
	// api.Get("/products", handlers.GetProducts)
	// api.Get("/products/:slug", handlers.GetProductDetail)

	// Admin routes
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware)

	// Test protected route
	admin.Get("/me", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"user_id": c.Locals("user_id"),
			"email":   c.Locals("email"),
		})
	})
    
    // Placeholder for routes to be implemented in future issues
    api.Get("/", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{"message": "Welcome to Machine Katalog API"})
    })
}
