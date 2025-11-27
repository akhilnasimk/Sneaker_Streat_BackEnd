package services

import (
	"errors"
	"fmt"

	"github.com/akhilnasimk/SS_backend/internal/dto"
	"github.com/akhilnasimk/SS_backend/internal/helpers"
	"github.com/akhilnasimk/SS_backend/internal/repositories/interfaces"
	"github.com/google/uuid"
)

type CartService interface {
	GetUserCartItems(userID uuid.UUID) (dto.CartResponse, error)
	AddItemToCart(userID uuid.UUID, productID uuid.UUID) error
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

func (s *cartService) AddItemToCart(userID uuid.UUID, productID uuid.UUID) error {

	//  Validate Product
	product, err := s.ProductRep.ProductById(productID)
	if err != nil {
		return errors.New("product not found")
	}

	// Optional: prevent adding unavailable products
	if product.StockCount <= 0 || product.DeletedAt.Valid {
		return errors.New("product not available")
	}

	//  Add to cart (auto-creates cart if missing)
	if err := s.cartRepo.AddItemToCart(userID, productID); err != nil {
		return fmt.Errorf("failed adding item to cart: %w", err)
	}

	return nil
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
