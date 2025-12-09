package dto

import (
	"time"

	"github.com/akhilnasimk/SS_backend/internal/models"
	"github.com/google/uuid"
)

type WishlistItemDTO struct {
	ID        uuid.UUID     `json:"id"`
	Product   ProductMinDTO `json:"product"`
	CreatedAt time.Time     `json:"created_at"`
}

type ProductMinDTO struct {
	ID          uuid.UUID         `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Price       int64             `json:"price"`
	Images      []ProductImageDTO `json:"images"`
}

type ProductImageDTO struct {
	ID       uuid.UUID `json:"id"`
	URL      string    `json:"url"`
	AltText  string    `json:"alt_text"`
	Priority int       `json:"priority"`
}

func ToWishlistDTO(items []models.Wishlist) []WishlistItemDTO {
	result := make([]WishlistItemDTO, 0)

	for _, w := range items {

		// pick only first image (if exists)
		var firstImg []ProductImageDTO
		if len(w.Product.Images) > 0 {
			img := w.Product.Images[0]
			firstImg = []ProductImageDTO{
				{
					ID:       img.ID,
					URL:      img.URL,
					AltText:  img.AltText,
					Priority: img.Priority,
				},
			}
		}

		result = append(result, WishlistItemDTO{
			ID:        w.ID,
			CreatedAt: w.CreatedAt,
			Product: ProductMinDTO{
				ID:          w.Product.ID,
				Name:        w.Product.Name,
				Description: w.Product.Description,
				Price:       w.Product.Price,
				Images:      firstImg, // ONLY FIRST IMAGE
			},
		})
	}

	return result
}

// Response structure
type WishlistStatusResponse struct {
	InWishlist bool      `json:"in_wishlist"`
	Exists     bool      `json:"exists"`
	Message    string    `json:"message"`
	WishlistID uuid.UUID `json:"wishlist_id,omitempty"`
	AddedAt    time.Time `json:"added_at,omitempty"`
}
