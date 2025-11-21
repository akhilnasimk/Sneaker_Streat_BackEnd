package models

import (
	"github.com/google/uuid"
)

type CartItem struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;index" json:"id"`
	CartID    uuid.UUID `gorm:"type:uuid;not null;index" json:"cart_id"`
	ProductID uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
	Product   *Product  `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"product"`

	Quantity int `gorm:"default:1" json:"quantity"`
}
