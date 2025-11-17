package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;index" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Price       int64     `gorm:"not null" json:"price"`
	StockCount  int       `gorm:"not null" json:"stock_count"`
	BaseImage   string    `json:"base_image"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`

	// 👇 Category Relation
	CategoryID uuid.UUID `gorm:"type:uuid;not null;index" json:"category_id"`
	Category   Category  `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE" json:"-"`

	// 👇 Images Relation
	Images []ProductImage `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"images"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
