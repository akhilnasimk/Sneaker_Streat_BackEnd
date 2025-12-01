package services

import (
	"fmt"
	"time"

	constent "github.com/akhilnasimk/SS_backend/internal/const"
	"github.com/akhilnasimk/SS_backend/internal/helpers"
	"github.com/akhilnasimk/SS_backend/internal/models"
	"github.com/akhilnasimk/SS_backend/internal/repositories/interfaces"
	"github.com/google/uuid"
)

type OrderService interface {
	GetAllOrders(userID string) ([]models.Order, error)
	CreateOrderFromCart(userIDString, shippingAddress, paymentMethod string) (*models.Order, error)
	CreateSingleOrder(userIDString string, productIDString string, quantity int, shippingAddress string, paymentMethod string) (*models.Order, error)
	CancelSingleOrderItem(orderItemIdString string) error
	CancelEntireOrder(orderIDStr string, userID uuid.UUID) error
	UpdateOrderStatus(orderID string, newStatus string) error
}

type orderService struct {
	OrderRepo interfaces.OrderRepository
	CartRepo  interfaces.CartRepository
}

func NewOrderService(orderRepo interfaces.OrderRepository, cartRepo interfaces.CartRepository) OrderService {
	return &orderService{
		OrderRepo: orderRepo,
		CartRepo:  cartRepo,
	}
}

const (
	CancellationWindowHours = 24 // Users can cancel within 24 hours
)

// -----------------------------------------------------------
// 1. Get all orders for user
// ----------------------------------------------------------
func (s *orderService) GetAllOrders(userID string) ([]models.Order, error) {
	U_id := helpers.StringToUUID(userID)
	return s.OrderRepo.FindAllOrders(U_id)
}

// -----------------------------------------------------------
// 2. Place Order From Entire Cart
// -----------------------------------------------------------
func (s *orderService) CreateOrderFromCart(userIDString, shippingAddress, paymentMethod string) (*models.Order, error) {

	userID := helpers.StringToUUID(userIDString)
	// 1. Fetch cart items
	cartItems, err := s.CartRepo.FindAllcartItemsOfUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed fetching cart items: %w", err)
	}
	if len(cartItems.CartItems) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	// 2. Calculate total amount
	var total float64
	for _, item := range cartItems.CartItems {
		total += float64(item.Quantity) * float64(item.Product.Price)
	}

	// 3. Create order model
	order := &models.Order{
		UserID:          userID,
		TotalAmount:     total,
		Status:          "pending",
		PaymentMethod:   paymentMethod,
		ShippingAddress: shippingAddress,
	}

	// 4. Create order + items together (handles stock, snapshots, cart deletion)
	if err := s.OrderRepo.CreateOrderWithItems(order, cartItems.CartItems); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}

// -----------------------------------------------------------
// 3. Place Order For A Single Product
// -----------------------------------------------------------
func (s *orderService) CreateSingleOrder(userIDString string, productIDString string, quantity int, shippingAddress string, paymentMethod string) (*models.Order, error) {
	// Validation
	if quantity <= 0 {
		return nil, fmt.Errorf("quantity must be greater than 0")
	}

	userID := helpers.StringToUUID(userIDString)
	productID := helpers.StringToUUID(productIDString)

	// Create order (TotalAmount will be set by repository after fetching product price)
	order := &models.Order{
		UserID:          userID,
		TotalAmount:     0, // Will be updated by repo
		Status:          "pending",
		PaymentMethod:   paymentMethod,
		ShippingAddress: shippingAddress,
	}

	// Create order with single item (handles stock, snapshot, total calculation)
	if err := s.OrderRepo.CreateSingleOrder(order, productID, quantity); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}

// cancel th singel orderitem
func (s *orderService) CancelSingleOrderItem(orderItemIdString string) error {
	// Parse ID
	id := helpers.StringToUUID(orderItemIdString)
	if id == uuid.Nil {
		return fmt.Errorf("invalid order item id")
	}

	// Fetch order item with parent order
	orderItem, err := s.OrderRepo.FindOrderItemByID(id)
	if err != nil {
		return fmt.Errorf("order item not found: %w", err)
	}

	//  Validate parent order only (2 checks)
	if err := s.validateParentOrder(orderItem.Order); err != nil {
		return err
	}

	// Execute cancellation
	return s.OrderRepo.CancelSingleOrderItem(orderItem.ID)
}

// CANCEL ENTIRE ORDER
func (s *orderService) CancelEntireOrder(orderIDStr string, userID uuid.UUID) error {
	// Parse ID
	orderID := helpers.StringToUUID(orderIDStr)
	if orderID == uuid.Nil {
		return fmt.Errorf("invalid order id")
	}

	// Fetch order
	order, err := s.OrderRepo.FindOrderByID(orderID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	// Check ownership
	if order.UserID != userID {
		return fmt.Errorf("unauthorized: this order does not belong to you")
	}

	//  Validate parent order only (2 checks)
	if err := s.validateParentOrder(order); err != nil {
		return err
	}

	// Execute cancellation
	return s.OrderRepo.CancelWholeOrder(order.ID)
}

// VALIDATION - Only 2 Checks on Parent Order bussinsess logic
func (s *orderService) validateParentOrder(order *models.Order) error {
	// Check 1: Status must be "pending"
	if order.Status != "pending" {
		return fmt.Errorf("cannot cancel - order status is '%s'. Only 'pending' orders can be cancelled", order.Status)
	}

	// Check 2: Must be within 24 hours
	if time.Since(order.CreatedAt) > CancellationWindowHours*time.Hour {
		hoursElapsed := int(time.Since(order.CreatedAt).Hours())
		return fmt.Errorf(
			"cancellation period expired - orders can only be cancelled within %d hours (order placed %d hours ago)",
			CancellationWindowHours,
			hoursElapsed,
		)
	}

	return nil
}

func (s *orderService) UpdateOrderStatus(orderID string, newStatus string) error {
	// validate enum
	if !constent.AllowedOrderStatus[newStatus] {
		return fmt.Errorf("invalid status value: %s", newStatus)
	}

	id, err := uuid.Parse(orderID)
	if err != nil {
		return fmt.Errorf("invalid order id")
	}

	// fetch current order
	order, err := s.OrderRepo.FindOrderByID(id)
	if err != nil {
		return err
	}

	current := order.Status

	// prevent delivered modifications
	if current == "delivered" {
		return fmt.Errorf("cannot modify delivered order")
	}

	// ensure correct transition
	nextAllowed, ok := constent.AllowedTransitions[current]
	if !ok || nextAllowed != newStatus {
		return fmt.Errorf("invalid transition %s → %s", current, newStatus)
	}

	// update
	return s.OrderRepo.UpdateOrderStatus(id, newStatus)
}
