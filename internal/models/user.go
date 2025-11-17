package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserName  string         `json:"username" gorm:"type:varchar(100);not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;type:varchar(100);not null"`
	Password  string         `json:"password" gorm:"not null"`
	Image     *string        `json:"image,omitempty"`
	Phone     *string        `json:"phone,omitempty"`
	Address   *string        `json:"address,omitempty"`
	IsAdmin   bool           `json:"is_admin" gorm:"default:false"`
	IsBlocked bool           `json:"is_blocked" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	UserRole  *string        `gorm:"default:'customer'"`

	// Relationships
	Cart      Cart       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"cart"`
	Orders    []Order    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"orders"`
	Wishlists []Wishlist `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"wishlists"`
}
