package handlers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/fianmhmmd/machine-catalog/backend/database"
	"github.com/fianmhmmd/machine-catalog/backend/models"
	"github.com/fianmhmmd/machine-catalog/backend/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/minio/minio-go/v7"
)

type ProductRequest struct {
	CategoryID     string       `json:"category_id" validate:"required"`
	Name           string       `json:"name" validate:"required"`
	Description    string       `json:"description"`
	Specifications models.JSONB `json:"specifications"`
	ContactPhone   string       `json:"contact_phone"`
	ContactName    string       `json:"contact_name"`
	IsPublished    bool         `json:"is_published"`
}

func GetProducts(c *fiber.Ctx) error {
	var products []models.Product
	
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit
	
	categorySlug := c.Query("category")
	search := c.Query("search")
	isAdmin := c.Locals("user_id") != nil

	query := database.DB.Preload("Category").Preload("Images")

	if !isAdmin {
		query = query.Where("is_published = ?", true)
	}

	if categorySlug != "" {
		var category models.Category
		if err := database.DB.Where("slug = ?", categorySlug).First(&category).Error; err == nil {
			query = query.Where("category_id = ?", category.ID)
		}
	}

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	var total int64
	query.Model(&models.Product{}).Count(&total)

	if err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&products).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not fetch products"})
	}

	return c.JSON(fiber.Map{
		"data": products,
		"meta": fiber.Map{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

func GetProductDetail(c *fiber.Ctx) error {
	slugStr := c.Params("slug")
	var product models.Product
	
	query := database.DB.Preload("Category").Preload("Images").Where("slug = ?", slugStr)
	
	if c.Locals("user_id") == nil {
		query = query.Where("is_published = ?", true)
	}

	if err := query.First(&product).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
	}

	return c.JSON(product)
}

func GetRelatedProducts(c *fiber.Ctx) error {
	slugStr := c.Params("slug")
	var product models.Product
	if err := database.DB.Where("slug = ?", slugStr).First(&product).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
	}

	var related []models.Product
	if err := database.DB.Preload("Category").Preload("Images").
		Where("category_id = ? AND id != ? AND is_published = ?", product.CategoryID, product.ID, true).
		Limit(4).Find(&related).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not fetch related products"})
	}

	return c.JSON(related)
}

func CreateProduct(c *fiber.Ctx) error {
	req := new(ProductRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	catID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid category ID"})
	}

	product := models.Product{
		CategoryID:     catID,
		Name:           req.Name,
		Slug:           slug.Make(req.Name),
		Description:    req.Description,
		Specifications: req.Specifications,
		ContactPhone:   req.ContactPhone,
		ContactName:    req.ContactName,
		IsPublished:    req.IsPublished,
	}

	if err := database.DB.Create(&product).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not create product"})
	}

	return c.Status(201).JSON(product)
}

func UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	req := new(ProductRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var product models.Product
	if err := database.DB.Where("id = ?", id).First(&product).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
	}

	catID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid category ID"})
	}

	product.CategoryID = catID
	product.Name = req.Name
	product.Slug = slug.Make(req.Name)
	product.Description = req.Description
	product.Specifications = req.Specifications
	product.ContactPhone = req.ContactPhone
	product.ContactName = req.ContactName
	product.IsPublished = req.IsPublished

	if err := database.DB.Save(&product).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not update product"})
	}

	return c.JSON(product)
}

func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var product models.Product
	if err := database.DB.Preload("Images").Where("id = ?", id).First(&product).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
	}

	bucketName := os.Getenv("MINIO_BUCKET")
	for _, img := range product.Images {
		objectName := filepath.Base(img.ImageURL)
		utils.MinioClient.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
	}

	if err := database.DB.Delete(&product).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not delete product"})
	}

	return c.JSON(fiber.Map{"message": "Product deleted successfully"})
}

func UploadProductImage(c *fiber.Ctx) error {
	productID := c.Params("id")
	
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Image file is required"})
	}

	if file.Size > 5*1024*1024 {
		return c.Status(400).JSON(fiber.Map{"error": "Image size exceeds 5MB"})
	}

	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		return c.Status(400).JSON(fiber.Map{"error": "Only jpg, png, and webp are allowed"})
	}

	objectName := fmt.Sprintf("product-%s-%d%s", productID, time.Now().Unix(), ext)
	bucketName := os.Getenv("MINIO_BUCKET")

	src, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not open file"})
	}
	defer src.Close()

	_, err = utils.MinioClient.PutObject(context.Background(), bucketName, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to upload to storage"})
	}

	pID, _ := uuid.Parse(productID)
	imageURL := fmt.Sprintf("/%s/%s", bucketName, objectName)
	
	var count int64
	database.DB.Model(&models.ProductImage{}).Where("product_id = ?", pID).Count(&count)
	
	productImage := models.ProductImage{
		ProductID: pID,
		ImageURL:  imageURL,
		IsPrimary: count == 0,
		SortOrder: int(count),
	}

	if err := database.DB.Create(&productImage).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save image reference to database"})
	}

	return c.JSON(productImage)
}

func DeleteProductImage(c *fiber.Ctx) error {
	imageID := c.Params("imageId")
	
	var img models.ProductImage
	if err := database.DB.Where("id = ?", imageID).First(&img).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Image not found"})
	}

	bucketName := os.Getenv("MINIO_BUCKET")
	objectName := filepath.Base(img.ImageURL)
	err := utils.MinioClient.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		fmt.Println("Warning: Failed to delete object from MinIO:", err)
	}

	if err := database.DB.Delete(&img).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not delete image from database"})
	}

	return c.JSON(fiber.Map{"message": "Image deleted successfully"})
}
