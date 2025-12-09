package services

import (
	"fmt"

	"github.com/akhilnasimk/SS_backend/internal/dto"
	"github.com/akhilnasimk/SS_backend/internal/helpers"
	"github.com/akhilnasimk/SS_backend/internal/repositories/interfaces"
	"github.com/google/uuid"
)

type CartService interface {
	GetUserCartItems(userID uuid.UUID) (dto.CartResponse, error)
	AddItemToCart(userID uuid.UUID, productID uuid.UUID) (*dto.CartItemResponse, error)
	IncOrDecCartItem(idstring string, oper string) error
	DeleteCartItem(idstring string) error
}

type cartService struct {
	cartRepo   interfaces.CartRepository
	ProductRep interfaces.ProductsRepository
}

// Constructor
func NewCartService(cartRepo interfaces.CartRepository, ProductRep interfaces.ProductsRepository) CartService {
	return &cartService{
		cartRepo:   cartRepo,
		ProductRep: ProductRep,
	}
}

// Fetch cart and map to DTO
func (s *cartService) GetUserCartItems(userID uuid.UUID) (dto.CartResponse, error) {
	cart, err := s.cartRepo.FindAllcartItemsOfUser(userID)
	if err != nil {
		return dto.CartResponse{}, err
	}

	return dto.MapCartToCartResponse(cart), nil
}
func (s *cartService) AddItemToCart(userID uuid.UUID, productID uuid.UUID) (*dto.CartItemResponse, error) {
	// Add to cart (no need for separate product validation as repo does it)
	cartItem, err := s.cartRepo.AddItemToCart(userID, productID)
	if err != nil {
		return nil, fmt.Errorf("failed adding item to cart: %w", err)
	}

	// Get product details (already preloaded in cartItem.Product)
	product := cartItem.Product

	// Build photo URL (handle empty images)
	photoURL := ""
	if len(product.Images) > 0 {
		photoURL = product.Images[0].URL
	}

	// Map to response
	resp := &dto.CartItemResponse{
		CartItemID:  cartItem.ID,
		ProductID:   product.ID,
		ProductName: product.Name,
		Price:       float64(product.Price),
		Photo:       photoURL,
		Total:       float64(product.Price) * float64(cartItem.Quantity),
		Quantity:    cartItem.Quantity,
		Catogory:    product.CategoryID,
	}

	return resp, nil
}

// patch the quantity of the cart items
func (s *cartService) IncOrDecCartItem(idstring string, oper string) error {
	if idstring == "" {
		return fmt.Errorf("the id is not given")
	}

	// Validate UUID
	id := helpers.StringToUUID(idstring)

	// Validate operation
	if oper != "inc" && oper != "dec" {
		return fmt.Errorf("invalid operation, must be 'inc' or 'dec'")
	}

	// Call repository
	if err := s.cartRepo.PatchQuantity(id, oper); err != nil {
		return err
	}

	return nil
}

func (s *cartService) DeleteCartItem(idstring string) error {
	id := helpers.StringToUUID(idstring)

	err := s.cartRepo.HardDeleteCartItem(id)
	if err != nil {
		return err
	}
	return nil
}
