package services

import (
	"errors"

	"github.com/akhilnasimk/SS_backend/internal/dto"
	"github.com/akhilnasimk/SS_backend/internal/helpers"
	"github.com/akhilnasimk/SS_backend/internal/models"
	"github.com/akhilnasimk/SS_backend/internal/repositories/interfaces"
	"gorm.io/gorm"
)

type WishlistService interface {
	GetAllWishlistItems(userID string) ([]dto.WishlistItemDTO, error)
	DeleteWishlistItem(userID, productID string) error
	ToggleWishlist(userID, productID string) (string, *models.Wishlist, error)
	CheckWishlistStatus(userIDStr, productIDStr string) (*dto.WishlistStatusResponse, error)
}

type wishlistService struct {
	wishlistRepo interfaces.WishlistRepository
	productRepo  interfaces.ProductsRepository
}

func NewWishlistService(wishlistRepo interfaces.WishlistRepository, productRepo interfaces.ProductsRepository) WishlistService {
	return &wishlistService{
		wishlistRepo: wishlistRepo,
		productRepo:  productRepo,
	}
}

// ---------- Get All Wishlist ----------
func (s *wishlistService) GetAllWishlistItems(userID string) ([]dto.WishlistItemDTO, error) {

	uid := helpers.StringToUUID(userID)

	// Fetch wishlist DB models
	items, err := s.wishlistRepo.FindAllWishItems(uid)
	if err != nil {
		return nil, errors.New("failed to fetch wishlist items")
	}

	// Convert DB model → DTO
	response := dto.ToWishlistDTO(items)

	return response, nil
}

// ---------- Delete Wishlist Item ----------
func (s *wishlistService) DeleteWishlistItem(userID, productID string) error {
	uid := helpers.StringToUUID(userID)
	pid := helpers.StringToUUID(productID)

	// Find the wishlist entry
	item, err := s.wishlistRepo.FindByUserAndProduct(uid, pid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("wishlist item not found")
		}
		return errors.New("failed to verify wishlist item")
	}

	// Delete it
	if err := s.wishlistRepo.DeleteWishlist(item.ID); err != nil {
		return errors.New("failed to delete wishlist item")
	}

	return nil
}

// ---------- Toggle Wishlist ----------
func (s *wishlistService) ToggleWishlist(userID, productID string) (string, *models.Wishlist, error) {

	uid := helpers.StringToUUID(userID)

	pid := helpers.StringToUUID(productID)

	// Verify product exists
	product, err := s.productRepo.FindById(pid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, errors.New("product not found")
		}
		return "", nil, errors.New("could not verify product")
	}

	if !product.IsActive {
		return "", nil, errors.New("product is not available")
	}

	// Toggle wishlist
	action, item, err := s.wishlistRepo.ToggleWishlist(uid, pid)
	if err != nil {
		return "", nil, err
	}

	return action, item, nil
}

// CheckWishlistStatus checks if a product is in user's wishlist
func (s *wishlistService) CheckWishlistStatus(userIDStr, productIDStr string) (*dto.WishlistStatusResponse, error) {

	userID := helpers.StringToUUID(userIDStr)
	productID := helpers.StringToUUID(productIDStr)
	wishlistItem, err := s.wishlistRepo.FindByUserAndProduct(userID, productID)
	if err != nil {
		// If record not found, product is not in wishlist
		if err == gorm.ErrRecordNotFound {
			return &dto.WishlistStatusResponse{
				InWishlist: false,
				Exists:     false,
				Message:    "Product not in wishlist",
			}, nil
		}
		// For other errors, return the error
		return nil, err
	}

	// Product found in wishlist
	return &dto.WishlistStatusResponse{
		InWishlist: true,
		Exists:     true,
		WishlistID: wishlistItem.ID,
		AddedAt:    wishlistItem.CreatedAt,
		Message:    "Product is in wishlist",
	}, nil
}
