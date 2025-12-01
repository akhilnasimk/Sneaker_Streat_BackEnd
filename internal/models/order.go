package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID              uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;index" json:"id"`
	UserID          uuid.UUID   `gorm:"type:uuid;not null;index" json:"user_id"`
	User            *User       `gorm:"foreignKey:UserID" json:"-"`
	TotalAmount     float64     `json:"total_amount"`
	Status          string      `gorm:"type:varchar(20);default:'pending'" json:"status"`
	PaymentMethod   string      `gorm:"type:varchar(20)" json:"payment_method"`
	ShippingAddress string      `gorm:"type:text" json:"shipping_address"`
	OrderItems      []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"order_items"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
	CancelledAt     *time.Time  `gorm:"default:NULL" json:"cancelled_at"`
}
