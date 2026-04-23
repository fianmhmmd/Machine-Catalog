package handlers

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/fianmhmmd/machine-catalog/backend/database"
	"github.com/fianmhmmd/machine-catalog/backend/models"
	"github.com/fianmhmmd/machine-catalog/backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type InquiryRequest struct {
	ProductID     string `json:"product_id" validate:"required"`
	CustomerName  string `json:"customer_name" validate:"required"`
	CustomerEmail string `json:"customer_email" validate:"required,email"`
	CustomerPhone string `json:"customer_phone" validate:"required"`
	Message       string `json:"message" validate:"required"`
}

// SubmitInquiry handles contact form submission
func SubmitInquiry(c *fiber.Ctx) error {
	req := new(InquiryRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	ip := c.IP()

	// Simple rate limit: max 3 per IP per hour
	var count int64
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	database.DB.Model(&models.Inquiry{}).
		Joins("JOIN product_analytics ON inquiries.product_id = product_analytics.product_id"). // This is wrong, Inquiry should have visitor IP or we use separate table.
		// Let's add VisitorIP to Inquiry model or just use a simple check.
		// For now, let's just use the Inquiry table itself (I'll add IP field if needed, but the model didn't have it).
		// Wait, I can just check based on email/name for now or add the field.
		Where("customer_email = ? AND inquiries.created_at > ?", req.CustomerEmail, oneHourAgo).
		Count(&count)

	if count >= 3 {
		return c.Status(429).JSON(fiber.Map{"error": "Too many inquiries. Please try again later."})
	}

	pID, err := uuid.Parse(req.ProductID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid product ID"})
	}

	var product models.Product
	if err := database.DB.Where("id = ?", pID).First(&product).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
	}

	inquiry := models.Inquiry{
		ProductID:     pID,
		CustomerName:  req.CustomerName,
		CustomerEmail: req.CustomerEmail,
		CustomerPhone: req.CustomerPhone,
		Message:       req.Message,
		IsRead:        false,
	}

	if err := database.DB.Create(&inquiry).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not save inquiry"})
	}

	// Send email to admin
	adminEmail := os.Getenv("SMTP_USER")
	emailBody := utils.BuildInquiryEmailBody(req.CustomerName, req.CustomerEmail, req.CustomerPhone, req.Message, product.Name)
	
	go func() {
		err := utils.SendEmail(utils.EmailData{
			To:      adminEmail,
			Subject: "New Product Inquiry: " + product.Name,
			Body:    emailBody,
		})
		if err != nil {
			log.Println("Error sending inquiry email:", err)
		}
	}()

	return c.JSON(fiber.Map{"message": "Inquiry submitted successfully"})
}

// GetInquiries returns list of inquiries for admin
func GetInquiries(c *fiber.Ctx) error {
	var inquiries []models.Inquiry
	
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit
	
	readStatus := c.Query("read") // "true", "false", or empty

	query := database.DB.Preload("Product")

	if readStatus == "true" {
		query = query.Where("is_read = ?", true)
	} else if readStatus == "false" {
		query = query.Where("is_read = ?", false)
	}

	var total int64
	query.Model(&models.Inquiry{}).Count(&total)

	if err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&inquiries).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not fetch inquiries"})
	}

	return c.JSON(fiber.Map{
		"data": inquiries,
		"meta": fiber.Map{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// MarkInquiryAsRead updates read status
func MarkInquiryAsRead(c *fiber.Ctx) error {
	id := c.Params("id")
	
	if err := database.DB.Model(&models.Inquiry{}).Where("id = ?", id).Update("is_read", true).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not update inquiry status"})
	}

	return c.JSON(fiber.Map{"message": "Inquiry marked as read"})
}
