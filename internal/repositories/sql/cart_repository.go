package sql

import (
	"errors"
	"fmt"

	"github.com/akhilnasimk/SS_backend/internal/models"
	"github.com/akhilnasimk/SS_backend/internal/repositories/interfaces"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type cartRepository struct {
	DB gorm.DB
}

func NewcartRepository(db gorm.DB) interfaces.CartRepository {
	return &cartRepository{
		DB: db,
	}
}

func (r *cartRepository) FindAllcartItemsOfUser(userID uuid.UUID) (models.Cart, error) {
	var cart models.Cart

	// Preload CartItems and Product for the given user
	err := r.DB.
		Preload("CartItems").
		Preload("CartItems.Product").
		Preload("CartItems.Product.Images").
		Where("user_id = ?", userID).
		First(&cart).Error

	if err != nil {
		return cart, err
	}

	return cart, nil
}

func (r *cartRepository) AddItemToCart(userID uuid.UUID, productID uuid.UUID) error {

	return r.DB.Transaction(func(tx *gorm.DB) error {

		//  Check if user already has a cart
		var cart models.Cart
		err := tx.Where("user_id = ?", userID).First(&cart).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No cart → create one
			newCart := models.Cart{
				ID:     uuid.New(),
				UserID: userID,
			}

			if err := tx.Create(&newCart).Error; err != nil {
				return fmt.Errorf("failed creating cart: %w", err)
			}

			cart = newCart
		} else if err != nil {
			return err
		}

		// Check if the cart already has this product
		var cartItem models.CartItem
		err = tx.Where("cart_id = ? AND product_id = ?", cart.ID, productID).First(&cartItem).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {

			// Insert new cart item
			newCartItem := models.CartItem{
				ID:        uuid.New(),
				CartID:    cart.ID,
				ProductID: productID,
				Quantity:  1,
			}

			if err := tx.Create(&newCartItem).Error; err != nil {
				return fmt.Errorf("failed creating cart item: %w", err)
			}

			return nil
		}

		if err != nil {
			return err
		}

		return nil
	})
}

func (r *cartRepository) PatchQuantity(id uuid.UUID, op string) error {
	switch op {
	case "inc":
		return r.DB.Model(&models.CartItem{}).
			Where("id = ?", id).
			Update("quantity", gorm.Expr("quantity + 1")).Error

	case "dec":
		// Prevent quantity from going below 1
		return r.DB.Model(&models.CartItem{}).
			Where("id = ?", id).
			Where("quantity > 1").
			Update("quantity", gorm.Expr("quantity - 1")).Error

	default:
		return errors.New("invalid operation")
	}
}

func (r *cartRepository) HardDeleteCartItem(id uuid.UUID) error {
    return r.DB.
        Where("id = ?", id).
        Delete(&models.CartItem{}).Error
}
