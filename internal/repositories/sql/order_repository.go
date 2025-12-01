package sql

import (
	"fmt"
	"time"

	"github.com/akhilnasimk/SS_backend/internal/models"
	"github.com/akhilnasimk/SS_backend/internal/repositories/interfaces"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type orderRepository struct {
	DB gorm.DB
}

func NewOrderRepository(db gorm.DB) interfaces.OrderRepository {
	return &orderRepository{
		DB: db,
	}
}

func (R *orderRepository) FindAllOrders(userID uuid.UUID) ([]models.Order, error) {
	var orders []models.Order

	// Preload order items, product, and product images
	err := R.DB.
		Preload("OrderItems", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC") // optional: order items by creation
		}).
		Preload("OrderItems.Product", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL AND is_active = ?", true)
		}).
		Preload("OrderItems.Product.Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL").Order("priority ASC")
		}).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error

	if err != nil {
		return []models.Order{}, err
	}

	// Always return array (even if empty)
	if len(orders) == 0 {
		return []models.Order{}, nil
	}

	return orders, nil
}

// methode for both single and cart ordering
func (r *orderRepository) CreateOrder(order *models.Order) error {
	return r.DB.Create(&order).Error
}

// ordering an entire cart
func (r *orderRepository) CreateOrderWithItems(order *models.Order, items []models.CartItem) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Create the order
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, ci := range items {
		// Fetch product with images
		var product models.Product
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Preload("Images", func(db *gorm.DB) *gorm.DB {
				return db.Order("priority ASC").Limit(1) // Get primary image
			}).
			Where("id = ?", ci.ProductID).
			First(&product).Error; err != nil {
			tx.Rollback()
			return err
		}

		// Check stock
		if product.StockCount < ci.Quantity {
			tx.Rollback()
			return fmt.Errorf("insufficient stock for product %s", product.Name)
		}

		// Decrease stock
		newStock := product.StockCount - ci.Quantity
		if err := tx.Model(&product).Update("stock_count", newStock).Error; err != nil {
			tx.Rollback()
			return err
		}

		// **CAPTURE PRODUCT SNAPSHOT**
		var productImage string
		if len(product.Images) > 0 {
			productImage = product.Images[0].URL // or ImageURL
		}

		// Create order item with snapshot
		orderItem := models.OrderItem{
			OrderID:      order.ID,
			ProductID:    ci.ProductID,
			ProductName:  product.Name,
			ProductImage: productImage,
			Quantity:     ci.Quantity,
			Price:        float64(product.Price),
			TotalPrice:   float64(ci.Quantity) * float64(product.Price),
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Delete cart items
	cartItemIDs := make([]uuid.UUID, len(items))
	for i, ci := range items {
		cartItemIDs[i] = ci.ID
	}

	if err := tx.Where("id IN ?", cartItemIDs).Delete(&models.CartItem{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// ordering a single item
func (r *orderRepository) CreateSingleOrder(order *models.Order, productID uuid.UUID, quantity int) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Fetch Product with images and lock
	var product models.Product
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("priority ASC").Limit(1)
		}).
		Where("id = ?", productID).
		First(&product).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Stock check
	if product.StockCount < quantity {
		tx.Rollback()
		return fmt.Errorf("insufficient stock for product %s", product.Name)
	}

	// **CALCULATE TOTAL**
	totalAmount := float64(quantity) * float64(product.Price)
	order.TotalAmount = totalAmount // ✅ SET THE TOTAL

	// Create Order
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Capture product snapshot
	var productImage string
	if len(product.Images) > 0 {
		productImage = product.Images[0].URL
	}

	// Create Order Item with snapshot
	orderItem := models.OrderItem{
		OrderID:      order.ID,
		ProductID:    productID,
		ProductName:  product.Name,
		ProductImage: productImage,
		Quantity:     quantity,
		Price:        float64(product.Price),
		TotalPrice:   totalAmount, // Same as order total for single item
	}

	if err := tx.Create(&orderItem).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Reduce stock
	newStock := product.StockCount - quantity
	if err := tx.Model(&product).Update("stock_count", newStock).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *orderRepository) CancelSingleOrderItem(orderItemID uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	now := time.Now()

	// 1. Fetch the order item
	var item models.OrderItem
	if err := tx.Where("id = ?", orderItemID).First(&item).Error; err != nil {
		tx.Rollback()
		return err
	}

	// already cancelled → stop
	if item.CancelledAt != nil {
		tx.Rollback()
		return fmt.Errorf("order item already cancelled")
	}

	// 2. Mark item cancelled
	if err := tx.Model(&item).Update("cancelled_at", now).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 3. Restore stock for THIS product
	var product models.Product
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", item.ProductID).
		First(&product).Error; err != nil {

		tx.Rollback()
		return err
	}

	restoredStock := product.StockCount + item.Quantity

	if err := tx.Model(&product).Update("stock_count", restoredStock).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 4. If all other items also cancelled → cancel whole order
	var activeCount int64
	tx.Model(&models.OrderItem{}).
		Where("order_id = ? AND cancelled_at IS NULL", item.OrderID).
		Count(&activeCount)

	if activeCount == 0 {
		// cancel order
		if err := tx.Model(&models.Order{}).
			Where("id = ?", item.OrderID).
			Update("status", "cancelled").Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *orderRepository) CancelWholeOrder(orderID uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	now := time.Now()

	// 1. Fetch all order items
	var items []models.OrderItem
	if err := tx.Where("order_id = ?", orderID).Find(&items).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 2. Restore stock for ALL items
	for _, item := range items {

		// skip already-cancelled items
		if item.CancelledAt != nil {
			continue
		}

		var product models.Product

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", item.ProductID).
			First(&product).Error; err != nil {

			tx.Rollback()
			return err
		}

		restored := product.StockCount + item.Quantity

		if err := tx.Model(&product).Update("stock_count", restored).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 3. Cancel ALL order items
	if err := tx.Model(&models.OrderItem{}).
		Where("order_id = ?", orderID).
		Update("cancelled_at", now).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 4. Cancel the order
	if err := tx.Model(&models.Order{}).
		Where("id = ?", orderID).
		Updates(map[string]interface{}{
			"status":       "cancelled",
			"cancelled_at": now,
		}).Error; err != nil {

		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *orderRepository) FindOrderItemByID(id uuid.UUID) (*models.OrderItem, error) {
	var orderItem models.OrderItem

	err := r.DB.
		Preload("Order"). // Load parent order for status checking
		Where("id = ?", id).
		First(&orderItem).Error

	if err != nil {
		return nil, err
	}

	return &orderItem, nil
}

func (r *orderRepository) FindOrderByID(id uuid.UUID) (*models.Order, error) {
	var order models.Order

	err := r.DB.
		Preload("OrderItems"). // ✅ Load all items for processing
		Where("id = ?", id).
		First(&order).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("order not found")
		}
		return nil, err
	}

	return &order, nil
}

func (r *orderRepository) UpdateOrderStatus(orderID uuid.UUID, newStatus string) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var order models.Order
	if err := tx.First(&order, "id = ?", orderID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// prevent update after delivered
	if order.Status == "delivered" {
		tx.Rollback()
		return fmt.Errorf("order already delivered, cannot update further")
	}

	// apply update
	if err := tx.Model(&order).Updates(map[string]interface{}{
		"status":     newStatus,
		"updated_at": time.Now(),
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
