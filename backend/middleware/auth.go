package middleware

import (
	"strings"

	"github.com/fianmhmmd/machine-catalog/backend/utils"
	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware protects routes by requiring a valid JWT access token
func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Authorization header missing"})
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid authorization format"})
	}

	tokenString := parts[1]
	claims, err := utils.VerifyToken(tokenString)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	// Set user info in locals for access in handlers
	c.Locals("user_id", claims["user_id"])
	c.Locals("email", claims["email"])

	return c.Next()
}
