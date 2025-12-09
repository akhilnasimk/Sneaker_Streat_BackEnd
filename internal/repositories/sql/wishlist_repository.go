package sql

import (
	"errors"
	"time"

	"github.com/akhilnasimk/SS_backend/internal/models"
	"github.com/akhilnasimk/SS_backend/internal/repositories/interfaces"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type wishlistRepository struct {
	DB *gorm.DB
}

func NewWishlistRepo(db *gorm.DB) interfaces.WishlistRepository {
	return &wishlistRepository{
		DB: db,
	}
}

// FindAllWishItems retrieves all wishlist items for a user
func (r *wishlistRepository) FindAllWishItems(userID uuid.UUID) ([]models.Wishlist, error) {
	var wishlist []models.Wishlist

	err := r.DB.
		Preload("Product", "deleted_at IS NULL AND is_active = ?", true).
		Preload("Product.Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL").Order("priority ASC")
		}).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&wishlist).Error

	if err != nil {
		return nil, err
	}

	// EMPTY RESULT → return empty array (NOT nil)
	if len(wishlist) == 0 {
		return []models.Wishlist{}, nil
	}

	return wishlist, nil
}

// DeleteWishlist removes a specific wishlist item by its ID
func (r *wishlistRepository) DeleteWishlist(wishlistID uuid.UUID) error {
	result := r.DB.Delete(&models.Wishlist{}, "id = ?", wishlistID)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// ToggleWishlist adds or removes a product from user's wishlist
func (r *wishlistRepository) ToggleWishlist(userID, productID uuid.UUID) (string, *models.Wishlist, error) {
	var existing models.Wishlist

	err := r.DB.
		Where("user_id = ? AND product_id = ?", userID, productID).
		First(&existing).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil, err
	}

	// If exists → remove it
	if err == nil {
		if delErr := r.DB.Delete(&existing).Error; delErr != nil {
			return "", nil, delErr
		}
		return "removed", nil, nil
	}

	// If not exists → add it
	newItem := models.Wishlist{
		UserID:    userID,
		ProductID: productID,
		CreatedAt: time.Now(),
	}

	if createErr := r.DB.Create(&newItem).Error; createErr != nil {
		return "", nil, createErr
	}

	// 🔥 PRELOAD DATA BEFORE RETURNING
	var fullItem models.Wishlist
	if preloadErr := r.DB.
		Preload("User").
		Preload("Product").
		Preload("Product.Images").
		Where("id = ?", newItem.ID).
		First(&fullItem).Error; preloadErr != nil {
		return "added", &newItem, nil // fallback
	}

	return "added", &fullItem, nil
}


func (r *wishlistRepository) FindByUserAndProduct(userID, productID uuid.UUID) (*models.Wishlist, error) {
	var item models.Wishlist
	err := r.DB.
		Where("user_id = ? AND product_id = ?", userID, productID).
		First(&item).Error

	if err != nil {
		return nil, err
	}
	return &item, nil
}
