package sql

import (
	"github.com/akhilnasimk/SS_backend/internal/models"
	"github.com/akhilnasimk/SS_backend/internal/repositories/interfaces"
	"gorm.io/gorm"
)

type productsRepository struct {
	DB gorm.DB
}

func NewProductsRepository(db gorm.DB) interfaces.ProductsRepository {
	return &productsRepository{
		DB: db,
	}
}

func (R *productsRepository) GetAllProducts(limit int, offset int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	// Count only active products
	if err := R.DB.Model(&models.Product{}).
		Where("is_active = ?", true).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated products
	result := R.DB.
		Preload("Images").
		Preload("Category").
		Where("is_active = ?", true).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&products)

	return products, total, result.Error
}
