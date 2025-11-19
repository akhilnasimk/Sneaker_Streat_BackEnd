package sql

import (
	"time"

	"github.com/akhilnasimk/SS_backend/internal/models"
	"github.com/akhilnasimk/SS_backend/internal/repositories/interfaces"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) interfaces.TokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) SaveRefreshToken(token *models.RefreshToken) error {
	return r.db.Create(token).Error
}

// Find a refresh token by its token string
func (r *tokenRepository) FindByToken(token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken

	// Search for token and ensure it's not revoked
	err := r.db.Where("token = ? AND revoked_at IS NULL", token).First(&refreshToken).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // no token found
		}
		return nil, err
	}

	return &refreshToken, nil
}

// Revoke a refresh token (mark as revoked)
func (r *tokenRepository) RevokeToken(token string) error {
	now := time.Now()

	// Update the revoked_at field
	return r.db.Model(&models.RefreshToken{}).
		Where("token = ?", token).
		Update("revoked_at", now).Error
}

func (r *tokenRepository) UpdateToken(userID uuid.UUID, newToken string, expiresAt time.Time) error {
	return r.db.Model(&models.RefreshToken{}).
		Where("user_id = ?", userID). // FIXED LINE
		Updates(map[string]interface{}{
			"token":      newToken,
			"expires_at": expiresAt,
			"revoked_at": nil,
		}).Error
}
