package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductImage struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ProductID uuid.UUID `gorm:"type:uuid;index"` // Foreign key
	URL       string    `json:"url"`
	AltText   string    `json:"alt_text"`

	// Optional: priority for ordering images
	Priority int `json:"priority"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
