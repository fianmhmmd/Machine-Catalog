package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) error {
	b.ID = uuid.New()
	return nil
}

type Admin struct {
	Base
	Email        string `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"not null" json:"-"`
	Name         string `gorm:"not null" json:"name"`
}

type Category struct {
	Base
	Name     string    `gorm:"not null" json:"name"`
	Slug     string    `gorm:"uniqueIndex;not null" json:"slug"`
	Products []Product `json:"products,omitempty"`
}

type Product struct {
	Base
	CategoryID     uuid.UUID      `gorm:"type:uuid;not null" json:"category_id"`
	Category       Category       `json:"category"`
	Name           string         `gorm:"not null" json:"name"`
	Slug           string         `gorm:"uniqueIndex;not null" json:"slug"`
	Description    string         `gorm:"type:text" json:"description"`
	Specifications JSONB          `gorm:"type:jsonb" json:"specifications"`
	ContactPhone   string         `json:"contact_phone"`
	ContactName    string         `json:"contact_name"`
	IsPublished    bool           `gorm:"default:false" json:"is_published"`
	Images         []ProductImage `json:"images"`
}

type ProductImage struct {
	Base
	ProductID uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
	ImageURL  string    `gorm:"not null" json:"image_url"`
	IsPrimary bool      `gorm:"default:false" json:"is_primary"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
}

type ProductAnalytics struct {
	Base
	ProductID uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
	EventType string    `gorm:"not null" json:"event_type"` // view, click
	VisitorIP string    `json:"visitor_ip"`
}

type Inquiry struct {
	Base
	ProductID     uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
	Product       Product   `json:"product"`
	CustomerName  string    `gorm:"not null" json:"customer_name"`
	CustomerEmail string    `gorm:"not null" json:"customer_email"`
	CustomerPhone string    `gorm:"not null" json:"customer_phone"`
	Message       string    `gorm:"type:text" json:"message"`
	IsRead        bool      `gorm:"default:false" json:"is_read"`
}

// Custom type for JSONB
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &j)
}
