package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/akhilnasimk/SS_backend/internal/dto"
	"github.com/akhilnasimk/SS_backend/internal/helpers"
	"github.com/akhilnasimk/SS_backend/internal/models"
	"github.com/akhilnasimk/SS_backend/internal/repositories/interfaces"
	"github.com/akhilnasimk/SS_backend/utils/cloudinary"
	"github.com/google/uuid"
)

type ProductsService interface {
	GetAllProducts(page, limit int, categoryID string, search string, minPrice, maxPrice int64) ([]models.Product, int64, error)
	GetProductById(idstring string) (dto.ProductResponse, error)
	CreateProduct(name, description string, price int64, stockCount int, categoryID uuid.UUID, files []*multipart.FileHeader) (models.Product, error)
}

type productsService struct {
	productRepo interfaces.ProductsRepository
}

func NewProductsService(repo interfaces.ProductsRepository) ProductsService {
	return &productsService{
		productRepo: repo,
	}
}

func (s *productsService) GetAllProducts(page, limit int, categoryID string, search string, minPrice, maxPrice int64) ([]models.Product, int64, error) {

	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	products, total, err := s.productRepo.GetAllProducts(limit, offset, categoryID, search, minPrice, maxPrice)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (s *productsService) GetProductById(idstring string) (dto.ProductResponse, error) {
	if idstring == "" {
		return dto.ProductResponse{}, fmt.Errorf("didn't send the id")
	}

	id := helpers.StringToUUID(idstring)
	product, err := s.productRepo.ProductById(id)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	// Map model to DTO
	return dto.ToProductResponse(product), nil
}

// the service became soo big so i changed the cloudinary entire service to another file in util 
func (s *productsService) CreateProduct(name, description string, price int64, stockCount int, categoryID uuid.UUID, files []*multipart.FileHeader) (models.Product, error) {
	// Set a reasonable timeout for the entire operation
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Prepare product model
	product := models.Product{
		Name:        name,
		Description: description,
		Price:       price,
		StockCount:  stockCount,
		CategoryID:  categoryID,
		IsActive:    true,
	}

	// Upload images to Cloudinary
	opts := cloudinary.DefaultUploadOptions("products")
	uploadResults, err := cloudinary.UploadMultiple(ctx, files, opts)
	if err != nil {
		return models.Product{}, fmt.Errorf("failed to upload images: %w", err)
	}

	// Convert upload results to ProductImage models
	var images []models.ProductImage
	for _, result := range uploadResults {
		images = append(images, models.ProductImage{
			URL:     result.URL,
			AltText: result.AltText,
		})
	}

	// Save to database
	createdProduct, err := s.productRepo.CreateProductWithImages(product, images)
	if err != nil {
		// Rollback: delete uploaded images from Cloudinary
		var uploadedURLs []string
		for _, img := range images {
			uploadedURLs = append(uploadedURLs, img.URL)
		}
		cloudinary.DeleteMultipleAsync(uploadedURLs)
		return models.Product{}, fmt.Errorf("failed to save product: %w", err)
	}

	return createdProduct, nil
}
