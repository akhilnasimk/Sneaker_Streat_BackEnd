package models

import (
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;index" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	User   *User     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`

	// Relationship to Cart Items
	CartItems []CartItem `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE" json:"items"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
