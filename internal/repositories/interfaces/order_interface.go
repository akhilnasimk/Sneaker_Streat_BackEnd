package interfaces

import (
	"github.com/akhilnasimk/SS_backend/internal/models"
	"github.com/google/uuid"
)

type OrderRepository interface {
	FindAllOrders(userID uuid.UUID) ([]models.Order, error)
	CreateOrderWithItems(order *models.Order, items []models.CartItem) error
	CreateSingleOrder(order *models.Order, productID uuid.UUID, quantity int) error
	CancelSingleOrderItem(orderItemID uuid.UUID) error
	CancelWholeOrder(orderID uuid.UUID) error
	FindOrderItemByID(id uuid.UUID) (*models.OrderItem, error)
	FindOrderByID(id uuid.UUID) (*models.Order, error)
	UpdateOrderStatus(orderID uuid.UUID, newStatus string) error
}
