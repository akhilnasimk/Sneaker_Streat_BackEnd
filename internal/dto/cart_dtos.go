package dto

import (
	"fmt"

	"github.com/akhilnasimk/SS_backend/internal/models"
	"github.com/google/uuid"
)

type CartItemResponse struct {
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	Price       float64   `json:"price"`
	Photo       string    `json:"photo"`
	Total       float64   `json:"total"`
	Quantity    int       `json:"quantity"`
}

type CartResponse struct {
	Items []CartItemResponse `json:"items"`
	Total float64            `json:"total"` // grand total of all cart items
}

// Mapping function helper for returning all cart items
func MapCartToCartResponse(cart models.Cart) CartResponse {
	var items []CartItemResponse
	var grandTotal float64

	for _, ci := range cart.CartItems {
		if ci.Product == nil || !ci.Product.IsActive {
			continue
		}

		// First non-deleted image
		var firstPhoto string
		for _, img := range ci.Product.Images {
			if !img.DeletedAt.Valid {
				firstPhoto = img.URL
				fmt.Println("the current image:", img)
				break
			}
		}

		fmt.Println("length of the image array is ", (ci.Product.Images))


		price := float64(ci.Product.Price)
		total := price * float64(ci.Quantity)
		grandTotal += total

		items = append(items, CartItemResponse{
			ProductID:   ci.Product.ID,
			ProductName: ci.Product.Name,
			Price:       price,
			Quantity:    ci.Quantity,
			Total:       total,
			Photo:       firstPhoto,
		})
	}

	return CartResponse{
		Items: items,
		Total: grandTotal,
	}
}

type AddCartRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required,min=1"`
}
