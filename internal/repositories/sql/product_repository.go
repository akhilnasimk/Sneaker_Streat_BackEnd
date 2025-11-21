package sql

import (
	"github.com/akhilnasimk/SS_backend/internal/models"
	"github.com/akhilnasimk/SS_backend/internal/repositories/interfaces"
	"github.com/google/uuid"
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

func (R *productsRepository) GetAllProducts(limit int, offset int, categoryID string, search string, minPrice int64, maxPrice int64) ([]models.Product, int64, error) {

	var products []models.Product
	var total int64

	db := R.DB.Model(&models.Product{}).Where("is_active = ?", true)

	// 🔥 Apply filters dynamically
	if categoryID != "" {
		db = db.Where("category_id = ?", categoryID)
	}

	if search != "" {
		db = db.Where("LOWER(name) LIKE LOWER(?)", "%"+search+"%")
	}

	if minPrice > 0 {
		db = db.Where("price >= ?", minPrice)
	}

	if maxPrice > 0 {
		db = db.Where("price <= ?", maxPrice)
	}

	// Count total with filters
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch results
	result := db.
		Preload("Images").
		Preload("Category").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&products)

	return products, total, result.Error
}

func (R *productsRepository) ProductById(id uuid.UUID) (models.Product, error) {
	var product models.Product

	err := R.DB.
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Preload("Category").
		First(&product, "id = ?", id).Error

	if err != nil {
		return models.Product{}, err
	}

	return product, nil
}

//admin prodduct managin

// products_repository.go

func (r *productsRepository) CreateProductWithImages(product models.Product, images []models.ProductImage) (models.Product, error) {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		// Create product
		if err := tx.Create(&product).Error; err != nil {
			return err
		}

		// Batch insert images in one query
		if len(images) > 0 {
			for i := range images {
				images[i].ProductID = product.ID
			}
			// This creates all images in a single INSERT statement
			if err := tx.Create(&images).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return models.Product{}, err
	}

	product.Images = images
	return product, nil
}