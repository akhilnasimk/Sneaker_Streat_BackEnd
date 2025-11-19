package interfaces

import "github.com/akhilnasimk/SS_backend/internal/models"

type ProductsRepository interface {
	GetAllProducts(limit int, offset int) ([]models.Product, int64, error)
}
