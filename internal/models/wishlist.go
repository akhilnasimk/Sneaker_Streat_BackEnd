package models

import (
	"time"

	"github.com/google/uuid"
)

type Wishlist struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index:idx_user_product,unique" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	ProductID uuid.UUID `gorm:"type:uuid;not null;index:idx_user_product,unique" json:"product_id"`
	Product   Product   `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"product"`
	CreatedAt time.Time `json:"created_at"`
}
