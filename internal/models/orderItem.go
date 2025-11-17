package models

import "github.com/google/uuid"

type OrderItem struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	OrderID    uuid.UUID `gorm:"type:uuid;not null;index" json:"order_id"`
	ProductID  uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
	Product    *Product  `gorm:"foreignKey:ProductID" json:"-"`
	Quantity   int       `json:"quantity"`
	Price      float64   `json:"price"`
	TotalPrice float64   `json:"total_price"`
}
