package models

import "github.com/google/uuid"

type ProductImage struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ProductID uuid.UUID `gorm:"type:uuid;index"`
	URL       string    `json:"url"`
	AltText   string    `json:"alt_text"`
	//preoiprity 
}
