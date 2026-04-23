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
	api.Get("/categories", handlers.GetCategories)
	api.Get("/products", handlers.GetProducts)
	api.Get("/products/:slug", handlers.GetProductDetail)
	api.Get("/products/:slug/related", handlers.GetRelatedProducts)

	// Admin routes
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware)

	// Admin Category routes
	admin.Post("/categories", handlers.CreateCategory)
	admin.Put("/categories/:id", handlers.UpdateCategory)
	admin.Delete("/categories/:id", handlers.DeleteCategory)

	// Admin Product routes
	admin.Get("/products", handlers.GetProducts)
	admin.Post("/products", handlers.CreateProduct)
	admin.Put("/products/:id", handlers.UpdateProduct)
	admin.Delete("/products/:id", handlers.DeleteProduct)

	// Admin Product Images
	admin.Post("/products/:id/images", handlers.UploadProductImage)
	admin.Delete("/products/:id/images/:imageId", handlers.DeleteProductImage)

	// Admin Analytics & Inquiry
	admin.Get("/analytics/overview", handlers.GetAnalyticsOverview)
	admin.Get("/inquiries", handlers.GetInquiries)
	admin.Put("/inquiries/:id/read", handlers.MarkInquiryAsRead)

	// Public Analytics & Inquiry
	api.Post("/products/:id/view", handlers.TrackView)
	api.Post("/products/:id/click", handlers.TrackClick)
	api.Post("/inquiry", handlers.SubmitInquiry)

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
