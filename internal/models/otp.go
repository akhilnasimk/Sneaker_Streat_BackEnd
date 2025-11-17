package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OTP struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID    *uuid.UUID     `gorm:"type:uuid;index" json:"user_id"` // nullable for signup flow
	User      *User          `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"user,omitempty"`
	Email     string         `gorm:"type:text;index:idx_email_purpose,unique" json:"email"` // explicitly store email, unique per purpose
	OTPCode   string         `gorm:"type:text;not null" json:"otp_code"`
	Purpose   string         `gorm:"type:text;not null;index:idx_email_purpose,unique" json:"purpose"` // login, signup, reset, order_verify
	ExpiresAt time.Time      `gorm:"not null" json:"expires_at"`
	IsUsed    bool           `gorm:"default:false" json:"is_used"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
