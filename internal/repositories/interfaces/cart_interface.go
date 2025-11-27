package interfaces

import (
	"github.com/akhilnasimk/SS_backend/internal/models"
	"github.com/google/uuid"
)

type CartRepository interface {
	FindAllcartItemsOfUser(userID uuid.UUID) (models.Cart, error)
	AddItemToCart(userID uuid.UUID, productID uuid.UUID) error
	PatchQuantity(id uuid.UUID, op string) error
	HardDeleteCartItem(id uuid.UUID) error
}
