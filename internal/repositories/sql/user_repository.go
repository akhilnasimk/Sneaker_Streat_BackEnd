package sql

import (
	"errors"
	"fmt"

	"github.com/akhilnasimk/SS_backend/internal/models"
	"github.com/akhilnasimk/SS_backend/internal/repositories/interfaces"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	DB gorm.DB
}

func NewUserReposetory(db gorm.DB) interfaces.UserRepository {
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

func (r *userRepository) GetAllUsersPaginated(limit, offset int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// Get total user count
	if err := r.DB.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated users, select only safe fields, order by creation date
	if err := r.DB.Select("id, user_name, email, created_at,user_role,is_admin,is_blocked").
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
