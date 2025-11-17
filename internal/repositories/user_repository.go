package repositories

import (
	"errors"
	"fmt"

	"github.com/akhilnasimk/SS_backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id uuid.UUID) (*models.User, error)
	PatchPasswordByEmail(email string, hashedPassword string) error
	PatchUser(id uuid.UUID, updates map[string]interface{}) error
}

type userRepository struct {
	DB gorm.DB
}

func NewUserReposetory(db gorm.DB) UserRepository {
	return &userRepository{
		DB: db,
	}
}

func (r *userRepository) CreateUser(user models.User) error {
	resp := r.DB.Create(&user)
	if resp.Error != nil {
		return resp.Error
	}
	return nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	resp := r.DB.Where("email = ?", email).First(&user)

	if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	resp := r.DB.First(&user, "id = ?", id)
	if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	return &user, nil
}

func (r *userRepository) PatchPasswordByEmail(email string, hashedPassword string) error {
	result := r.DB.Model(&models.User{}).
		Where("email = ?", email).
		Update("password", hashedPassword)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepository) PatchUser(id uuid.UUID, updates map[string]interface{}) error {
	result := r.DB.Model(&models.User{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
