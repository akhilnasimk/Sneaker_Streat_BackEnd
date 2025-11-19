package interfaces

import (
	"github.com/akhilnasimk/SS_backend/internal/models"
	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(user models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id uuid.UUID) (*models.User, error)
	PatchPasswordByEmail(email string, hashedPassword string) error
	PatchUser(id uuid.UUID, updates map[string]interface{}) error
	GetAllUsersPaginated(limit, offset int) ([]models.User, int64, error)
}
