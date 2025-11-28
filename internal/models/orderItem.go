package models

import "github.com/google/uuid"

type OrderItem struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`

	// Order relation
	OrderID uuid.UUID `gorm:"type:uuid;not null;index" json:"order_id"`
	Order   *Order    `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE;" json:"-"`

	// Product relation
	ProductID uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
	Product   *Product  `gorm:"foreignKey:ProductID;constraint:OnDelete:SET NULL;" json:"-"`

	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
	TotalPrice float64 `json:"total_price"`
}

