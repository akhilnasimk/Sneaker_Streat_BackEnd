package services

import (
	"context"
	"fmt"
	"log"
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
	GetAllProducts(page, limit int, categoryID string, search string, minPrice, maxPrice int64, userRole string) ([]models.Product, int64, error)
	GetProductById(idstring string) (dto.ProductResponse, error)
	CreateProduct(name, description string, price int64, stockCount int, categoryID uuid.UUID, files []*multipart.FileHeader) (models.Product, error)
	GetAllCategory() ([]dto.CategoryResponse, error)
	UpdateProduct(id uuid.UUID, req dto.UpdateProductRequest) error
	ToggleProductAvailability(idString string) error
	DeleteProduct(idString string) error
}

type productsService struct {
	productRepo interfaces.ProductsRepository
}

func NewProductsService(repo interfaces.ProductsRepository) ProductsService {
	return &productsService{
		productRepo: repo,
	}
}

func (s *productsService) GetAllProducts(page, limit int, categoryID string, search string, minPrice, maxPrice int64, userRole string) ([]models.Product, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	// Admin sees everything (including deleted), others see only active non-deleted
	includeDeleted := userRole == "admin"

	products, total, err := s.productRepo.GetAllProducts(limit, offset, categoryID, search, minPrice, maxPrice, includeDeleted)
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

func (s *productsService) GetAllCategory() ([]dto.CategoryResponse, error) {
	// Fetch categories from repository
	categories, err := s.productRepo.FindAllCategory()
	if err != nil {
		return nil, err
	}

	// Map models → DTO
	var resp []dto.CategoryResponse
	resp = make([]dto.CategoryResponse, 0, len(categories)) // allocate properly

	for _, c := range categories {
		resp = append(resp, dto.CategoryResponse{
			ID:        c.ID,
			Name:      c.Name,
			CreatedAt: c.CreatedAt,
		})
	}

	return resp, nil
}

func (s *productsService) UpdateProduct(id uuid.UUID, req dto.UpdateProductRequest) error {
	//  Fetch existing product
	product, err := s.productRepo.FindById(id)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	// Update basic fields
	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.StockCount = req.StockCount

	catID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return fmt.Errorf("invalid category ID: %w", err)
	}
	product.CategoryID = catID

	// Handle Image Deletion
	// Only process if user explicitly sent keep_images[] (even if empty)
	var removedURLs []string
	fmt.Println("update product service  p:id:", product.ID)
	fmt.Println("update product service  keep image :", req.KeepImages)
	if req.KeepImages != nil {
		removedURLs, err = s.productRepo.DeleteImagesNotIn(product.ID, req.KeepImages)
		if err != nil {
			return fmt.Errorf("failed to delete old images: %w", err)
		}

		// Remove deleted images from the in-memory slice so Save won't re-insert them
		if len(req.KeepImages) == 0 {
			// empty slice => delete all images
			product.Images = []models.ProductImage{}
		} else {
			keepSet := make(map[string]struct{}, len(req.KeepImages))
			for _, u := range req.KeepImages {
				keepSet[u] = struct{}{}
			}
			filtered := make([]models.ProductImage, 0, len(product.Images))
			for _, img := range product.Images {
				if _, ok := keepSet[img.URL]; ok {
					filtered = append(filtered, img)
				}
			}
			product.Images = filtered
		}

		// Delete from Cloudinary (async, fire-and-forget)
		if len(removedURLs) > 0 {
			log.Printf("Deleting %d images from Cloudinary for product %s", len(removedURLs), product.ID)
			cloudinary.DeleteMultipleAsync(removedURLs)
		}
	}

	// Upload New Images
	if len(req.NewImages) > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		results, err := cloudinary.UploadMultiple(
			ctx,
			req.NewImages,
			cloudinary.DefaultUploadOptions("products"),
		)

		if err != nil {
			return fmt.Errorf("failed to upload new images: %w", err)
		}

		// Add successfully uploaded images to product
		for _, img := range results {
			if img.Error == nil {
				product.Images = append(product.Images, models.ProductImage{
					URL:       img.URL,
					AltText:   img.AltText,
					ProductID: product.ID,
				})
			} else {
				log.Printf("Failed to upload image %s: %v", img.AltText, img.Error)
			}
		}
	}

	// Save everything to database
	if err := s.productRepo.UpdateProduct(product); err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	return nil
}

func (s *productsService) ToggleProductAvailability(idString string) error {
	id, err := uuid.Parse(idString)
	if err != nil {
		return fmt.Errorf("invalid product ID: %w", err)
	}

	if err := s.productRepo.ToggleActive(id); err != nil {
		return err
	}

	return nil
}

func (s *productsService) DeleteProduct(idString string) error {
	id, err := uuid.Parse(idString)
	if err != nil {
		return fmt.Errorf("invalid product ID: %w", err)
	}

	if err := s.productRepo.DeleteProduct(id); err != nil {
		return err
	}

	return nil
}
