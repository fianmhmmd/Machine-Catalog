package database

import (
	"fmt"
	"log"
	"os"

	"github.com/fianmhmmd/machine-catalog/backend/models"
	"github.com/fianmhmmd/machine-catalog/backend/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected")

	// Auto migrations
	db.AutoMigrate(
		&models.Admin{},
		&models.Category{},
		&models.Product{},
		&models.ProductImage{},
		&models.ProductAnalytics{},
		&models.Inquiry{},
	)

	seedAdmin(db)

	DB = db
}

func seedAdmin(db *gorm.DB) {
	var count int64
	db.Model(&models.Admin{}).Count(&count)
	if count == 0 {
		hashedPassword, _ := utils.HashPassword("password")
		admin := models.Admin{
			Name:         "Admin",
			Email:        "admin@example.com",
			PasswordHash: hashedPassword,
		}
		db.Create(&admin)
		log.Println("Default admin seeded: admin@example.com / password")
	}
}
