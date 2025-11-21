package dto

import (
	"time"

	"github.com/akhilnasimk/SS_backend/internal/models"
	"github.com/google/uuid"
)

type ProductImageResponse struct {
	ID      uuid.UUID `json:"id"`
	URL     string    `json:"url"`
	AltText string    `json:"alt_text"`
}

type ProductResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	StockCount  int       `json:"stock_count"`
	IsActive    bool      `json:"is_active"`

	CategoryID uuid.UUID `json:"category_id"`
	Category   string    `json:"category_name"` // optional, if you preload category

	Images []ProductImageResponse `json:"images"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func ToProductResponse(p models.Product) ProductResponse {
	images := make([]ProductImageResponse, len(p.Images))
	for i, img := range p.Images {
		images[i] = ProductImageResponse{
			ID:      img.ID,
			URL:     img.URL,
			AltText: img.AltText,
		}
	}

	return ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		StockCount:  p.StockCount,
		IsActive:    p.IsActive,
		CategoryID:  p.CategoryID,
		Category:    p.Category.Name,
		Images:      images,
		CreatedAt:   p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   p.UpdatedAt.Format(time.RFC3339),
	}
}
