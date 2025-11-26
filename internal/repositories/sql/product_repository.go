package sql

import (
	"fmt"

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

func (r *productsRepository) GetAllProducts(limit int, offset int, categoryID string, search string, minPrice int64, maxPrice int64, includeDeleted bool) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	db := r.DB.Model(&models.Product{})

	// Show deleted products for admin
	if includeDeleted {
		db = db.Unscoped() // Include soft-deleted records
	}

	// Filter by active status (non-admin only)
	if !includeDeleted {
		db = db.Where("is_active = ?", true)
	}

	// Apply filters
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

// catogories to map on the front end
func (r *productsRepository) FindAllCategory() ([]models.Category, error) {
	var categories []models.Category
	resp := r.DB.Find(&categories)
	if resp.Error != nil {
		return categories, resp.Error
	}
	return categories, nil
}

// find by id for helper to anothermethode in hte repo
func (r *productsRepository) FindById(id uuid.UUID) (*models.Product, error) {
	var product models.Product
	// Preload existing images
	err := r.DB.Preload("Images").Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// product upadation
func (r *productsRepository) UpdateProduct(product *models.Product) error {
	// Use Session to ensure associations are saved
	return r.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(product).Error
}

// delete product and related images that is not needed
func (r *productsRepository) DeleteImagesNotIn(productID uuid.UUID, urlsToKeep []string) ([]string, error) {
	var toDelete []models.ProductImage

	// Build query
	query := r.DB.Where("product_id = ?", productID)
	if len(urlsToKeep) > 0 {
		query = query.Where("url NOT IN ?", urlsToKeep)
	}

	// Fetch only the fields we need (ID and URL) for better performance
	if err := query.Select("id", "url").Find(&toDelete).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch images: %w", err)
	}

	if len(toDelete) == 0 {
		return []string{}, nil
	}

	// Extract IDs and URLs in a single loop
	ids := make([]uuid.UUID, len(toDelete))
	urls := make([]string, len(toDelete))
	for i, img := range toDelete {
		ids[i] = img.ID
		urls[i] = img.URL
	}

	// Delete images
	if err := r.DB.Where("id IN ?", ids).Delete(&models.ProductImage{}).Error; err != nil {
		return nil, fmt.Errorf("failed to delete images: %w", err)
	}

	return urls, nil
}

func (r *productsRepository) ToggleActive(id uuid.UUID) error {
	result := r.DB.Model(&models.Product{}).
		Where("id = ?", id).
		Update("is_active", gorm.Expr("NOT is_active"))

	if result.Error != nil {
		return fmt.Errorf("failed to toggle product active status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("product not found with id: %s", id)
	}

	return nil
}

func (r *productsRepository) DeleteProduct(id uuid.UUID) error {
	result := r.DB.Delete(&models.Product{}, id)

	if result.Error != nil {
		return fmt.Errorf("failed to delete product: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("product not found with id: %s", id)
	}

	return nil
}
