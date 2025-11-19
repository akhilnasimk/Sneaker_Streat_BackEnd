package interfaces

import (
	"time"

	"github.com/akhilnasimk/SS_backend/internal/models"
	"github.com/google/uuid"
)

type TokenRepository interface {
	SaveRefreshToken(token *models.RefreshToken) error
	FindByToken(token string) (*models.RefreshToken, error)
	RevokeToken(token string) error
	UpdateToken(id uuid.UUID, newToken string, expiresAt time.Time) error
}
