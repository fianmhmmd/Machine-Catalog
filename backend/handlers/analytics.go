package handlers

import (
	"time"

	"github.com/fianmhmmd/machine-catalog/backend/database"
	"github.com/fianmhmmd/machine-catalog/backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// TrackView records a product view
func TrackView(c *fiber.Ctx) error {
	return trackEvent(c, "view")
}

// TrackClick records a product contact click
func TrackClick(c *fiber.Ctx) error {
	return trackEvent(c, "click")
}

func trackEvent(c *fiber.Ctx, eventType string) error {
	productID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid product ID"})
	}

	ip := c.IP()

	// Simple deduplication: 1 event per IP per product per 24 hours
	var existing models.ProductAnalytics
	oneDayAgo := time.Now().Add(-24 * time.Hour)
	err = database.DB.Where("product_id = ? AND event_type = ? AND visitor_ip = ? AND created_at > ?", 
		productID, eventType, ip, oneDayAgo).First(&existing).Error

	if err == nil {
		// Event already recorded within 24h, skip
		return c.JSON(fiber.Map{"status": "skipped", "message": "Event already recorded"})
	}

	analytics := models.ProductAnalytics{
		ProductID: productID,
		EventType: eventType,
		VisitorIP: ip,
	}

	if err := database.DB.Create(&analytics).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not record analytics"})
	}

	return c.JSON(fiber.Map{"status": "success"})
}

// GetAnalyticsOverview returns summary data for admin
func GetAnalyticsOverview(c *fiber.Ctx) error {
	sevenDaysAgo := time.Now().Add(-7 * 24 * time.Hour)

	var totalViews int64
	database.DB.Model(&models.ProductAnalytics{}).Where("event_type = ? AND created_at > ?", "view", sevenDaysAgo).Count(&totalViews)

	var totalClicks int64
	database.DB.Model(&models.ProductAnalytics{}).Where("event_type = ? AND created_at > ?", "click", sevenDaysAgo).Count(&totalClicks)

	// Top 5 products by views
	type TopProduct struct {
		ID    uuid.UUID `json:"id"`
		Name  string    `json:"name"`
		Views int64     `json:"views"`
	}
	var topProducts []TopProduct
	database.DB.Table("products").
		Select("products.id, products.name, count(product_analytics.id) as views").
		Joins("left join product_analytics on products.id = product_analytics.product_id").
		Where("product_analytics.event_type = ? AND product_analytics.created_at > ?", "view", sevenDaysAgo).
		Group("products.id").
		Order("views DESC").
		Limit(5).
		Scan(&topProducts)

	return c.JSON(fiber.Map{
		"overview": fiber.Map{
			"total_views":  totalViews,
			"total_clicks": totalClicks,
			"period_days":  7,
		},
		"top_products": topProducts,
	})
}
