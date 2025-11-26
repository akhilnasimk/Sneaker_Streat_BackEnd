package helpers

import (
	"fmt"

	"github.com/akhilnasimk/SS_backend/internal/dto"
)

// ValidateUpdateProductRequest validates the update product request
func ValidateUpdateProductRequest(req dto.UpdateProductRequest) error {
	if req.Name == "" {
		return fmt.Errorf("product name is required")
	}
	if req.Price <= 0 {
		return fmt.Errorf("price must be greater than 0")
	}
	if req.StockCount < 0 {
		return fmt.Errorf("stock count cannot be negative")
	}
	if req.CategoryID == "" {
		return fmt.Errorf("category ID is required")
	}
	return nil
}

// ValidateCreateProductRequest validates the create product request
// func ValidateCreateProductRequest(req dto.CreateProductRequest) error {
// 	if req.Name == "" {
// 		return fmt.Errorf("product name is required")
// 	}
// 	if req.Price <= 0 {
// 		return fmt.Errorf("price must be greater than 0")
// 	}
// 	if req.StockCount < 0 {
// 		return fmt.Errorf("stock count cannot be negative")
// 	}
// 	if req.CategoryID == "" {
// 		return fmt.Errorf("category ID is required")
// 	}
// 	if len(req.Images) == 0 {
// 		return fmt.Errorf("at least one product image is required")
// 	}
// 	return nil
// }
