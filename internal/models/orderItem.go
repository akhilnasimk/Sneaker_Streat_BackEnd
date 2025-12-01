package models

import (
	"time"

	"github.com/google/uuid"
)
type OrderItem struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`

	// Order relation
	OrderID uuid.UUID `gorm:"type:uuid;not null;index" json:"order_id"`
	Order   *Order    `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE;" json:"-"`

	// Product reference (may be NULL after deletion)
	ProductID uuid.UUID `gorm:"type:uuid;index" json:"product_id"`
	Product   *Product  `gorm:"foreignKey:ProductID;constraint:OnDelete:SET NULL;" json:"-"`

	// SNAPSHOT FIELDS
	ProductName  string `gorm:"not null" json:"product_name"`
	ProductImage string `json:"product_image"`

	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
	TotalPrice float64 `json:"total_price"`

	// CANCELLATION SUPPORT
	CancelledAt *time.Time `gorm:"default:NULL" json:"cancelled_at"`

	// ➕ ADDED (Safe timestamps)
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
