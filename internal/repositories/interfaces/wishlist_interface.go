package interfaces

import (
	"github.com/akhilnasimk/SS_backend/internal/models"
	"github.com/google/uuid"
)

type WishlistRepository interface {
	FindAllWishItems(id uuid.UUID) ([]models.Wishlist, error)
	DeleteWishlist(id uuid.UUID) error
	ToggleWishlist(userID, productID uuid.UUID) (string, *models.Wishlist, error)
	FindByUserAndProduct(userID, productID uuid.UUID) (*models.Wishlist, error)
}
