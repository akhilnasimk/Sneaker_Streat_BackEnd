package models

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	Token     string     `gorm:"type:text;not null" json:"token"` // store the actual refresh token
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`      // token expiry
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at"` // optional: when token was revoked

	User User `gorm:"foreignKey:UserID"` // relationship to User
}
