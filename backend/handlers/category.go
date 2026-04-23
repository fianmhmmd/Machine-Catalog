package handlers

import (
	"github.com/fianmhmmd/machine-catalog/backend/database"
	"github.com/fianmhmmd/machine-catalog/backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
)

type CategoryRequest struct {
	Name string `json:"name" validate:"required"`
}

// GetCategories returns all categories
func GetCategories(c *fiber.Ctx) error {
	var categories []models.Category
	if err := database.DB.Find(&categories).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not fetch categories"})
	}
	return c.JSON(categories)
}

// CreateCategory handles category creation
func CreateCategory(c *fiber.Ctx) error {
	req := new(CategoryRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	category := models.Category{
		Name: req.Name,
		Slug: slug.Make(req.Name),
	}

	if err := database.DB.Create(&category).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create category"})
	}

	return c.Status(201).JSON(category)
}

// UpdateCategory handles category updates
func UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	req := new(CategoryRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var category models.Category
	if err := database.DB.Where("id = ?", id).First(&category).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Category not found"})
	}

	category.Name = req.Name
	category.Slug = slug.Make(req.Name)

	if err := database.DB.Save(&category).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not update category"})
	}

	return c.JSON(category)
}

// DeleteCategory handles category deletion
func DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	
	var category models.Category
	if err := database.DB.Where("id = ?", id).First(&category).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Category not found"})
	}

	var count int64
	database.DB.Model(&models.Product{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Category is being used by products and cannot be deleted"})
	}

	if err := database.DB.Delete(&category).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not delete category"})
	}

	return c.JSON(fiber.Map{"message": "Category deleted successfully"})
}
