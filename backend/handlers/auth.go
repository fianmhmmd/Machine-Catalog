package handlers

import (
	"github.com/fianmhmmd/machine-catalog/backend/database"
	"github.com/fianmhmmd/machine-catalog/backend/models"
	"github.com/fianmhmmd/machine-catalog/backend/utils"
	"github.com/gofiber/fiber/v2"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// Login handles admin login
func Login(c *fiber.Ctx) error {
	req := new(LoginRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var admin models.Admin
	if err := database.DB.Where("email = ?", req.Email).First(&admin).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	if !utils.CheckPasswordHash(req.Password, admin.PasswordHash) {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	at, rt, err := utils.GenerateToken(utils.TokenPayload{
		UserID: admin.ID.String(),
		Email:  admin.Email,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not generate tokens"})
	}

	return c.JSON(fiber.Map{
		"access_token":  at,
		"refresh_token": rt,
		"admin": fiber.Map{
			"id":    admin.ID,
			"name":  admin.Name,
			"email": admin.Email,
		},
	})
}

// Refresh handles token refresh
func Refresh(c *fiber.Ctx) error {
	req := new(RefreshRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	claims, err := utils.VerifyToken(req.RefreshToken)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid or expired refresh token"})
	}

	userID := claims["user_id"].(string)
	var admin models.Admin
	if err := database.DB.Where("id = ?", userID).First(&admin).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Admin not found"})
	}

	at, rt, err := utils.GenerateToken(utils.TokenPayload{
		UserID: admin.ID.String(),
		Email:  admin.Email,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not generate tokens"})
	}

	return c.JSON(fiber.Map{
		"access_token":  at,
		"refresh_token": rt,
	})
}
