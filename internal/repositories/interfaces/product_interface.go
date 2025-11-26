package interfaces

import (
	"github.com/akhilnasimk/SS_backend/internal/models"
	"github.com/google/uuid"
)

type ProductsRepository interface {
	GetAllProducts(limit int, offset int, categoryID string, search string, minPrice int64, maxPrice int64) ([]models.Product, int64, error)
	ProductById(id uuid.UUID) (models.Product, error)
	CreateProductWithImages(product models.Product, images []models.ProductImage) (models.Product, error)
	FindAllCategory() ([]models.Category, error)
	UpdateProduct(product *models.Product) error
	DeleteImagesNotIn(productID uuid.UUID, urlsToKeep []string) ([]string, error)
	FindById(id uuid.UUID) (*models.Product, error)
}
