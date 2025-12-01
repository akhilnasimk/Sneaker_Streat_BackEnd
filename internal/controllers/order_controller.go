package controllers

import (
	"net/http"

	"github.com/akhilnasimk/SS_backend/internal/dto"
	"github.com/akhilnasimk/SS_backend/internal/helpers"
	"github.com/akhilnasimk/SS_backend/internal/services"
	"github.com/akhilnasimk/SS_backend/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderController struct {
	OrderService services.OrderService
}

func NewOrderController(service services.OrderService) OrderController {
	return OrderController{
		OrderService: service,
	}
}

func (C *OrderController) GetAllOrders(ctx *gin.Context) {
	id, exist := ctx.Get("UserID")
	if !exist {
		ctx.JSON(400, response.Failure("id not availabel", nil))
	}

	orders, err := C.OrderService.GetAllOrders(id.(string))

	if err != nil {
		ctx.JSON(400, response.Failure("Failed to get the orders ", err))
		return
	}

	ctx.JSON(200, response.Success("order fetch sucess", orders))
}

func (C *OrderController) AddSingleItemOrder(ctx *gin.Context) {
	// Get User ID
	id, exist := ctx.Get("UserID")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, response.Failure("user id not available", nil))
		return
	}

	userID := id.(string)

	// Get Product ID
	productID := ctx.Param("product_id")
	if productID == "" {
		ctx.JSON(http.StatusBadRequest, response.Failure("product id not available", nil))
		return
	}

	// Bind Body
	var orderReq dto.CreateSingleOrderDTO
	if err := ctx.ShouldBindJSON(&orderReq); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Failure("failed to bind request body", err.Error()))
		return
	}

	// Service Call
	orderRes, err := C.OrderService.CreateSingleOrder(
		userID,
		productID,
		orderReq.Quantity,
		orderReq.ShippingAddress,
		orderReq.PaymentMethod,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Failure("failed to place order", err.Error()))
		return
	}

	// SUCCESS → 201 CREATED
	ctx.JSON(http.StatusCreated, response.Success("order placed successfully", orderRes))
}

func (C *OrderController) AddCartOrder(ctx *gin.Context) {
	userID, exists := ctx.Get("UserID")
	if !exists {
		ctx.JSON(400, response.Failure("user id missing", nil))
		return
	}

	var orderReq dto.CreateCartOrderDTO
	if err := ctx.ShouldBindJSON(&orderReq); err != nil {
		ctx.JSON(400, response.Failure("invalid request body", err.Error()))
		return
	}

	orderRes, err := C.OrderService.CreateOrderFromCart(
		userID.(string),
		orderReq.ShippingAddress,
		orderReq.PaymentMethod,
	)

	if err != nil {
		ctx.JSON(500, response.Failure("failed to create order from cart", err.Error()))
		return
	}

	ctx.JSON(201, response.Success("order created from cart", orderRes))
}

// CancelOrderItem - DELETE /orders/items/:item_id/cancel
func (c *OrderController) CancelOrderItem(ctx *gin.Context) {
	itemID := ctx.Param("item_id")

	// Call service (service validates the ID)
	if err := c.OrderService.CancelSingleOrderItem(itemID); err != nil {
		// Check error type for appropriate status code
		ctx.JSON(http.StatusBadRequest, response.Failure(err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, response.Success("Order item cancelled successfully", nil))
}

// CancelOrder - DELETE /orders/:order_id/cancel
func (c *OrderController) CancelOrder(ctx *gin.Context) {
	orderID := ctx.Param("order_id")

	// Read UserID stored by middleware
	userIDRaw, exists := ctx.Get("UserID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.Failure("Unauthorized", nil))
		return
	}

	// Safe type check
	userIDStr, ok := userIDRaw.(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, response.Failure("Invalid UserID format", nil))
		return
	}

	// Convert string → UUID
	id := helpers.StringToUUID(userIDStr)
	if id == uuid.Nil {
		ctx.JSON(http.StatusBadRequest, response.Failure("Invalid user UUID", nil))
		return
	}

	// Call service
	if err := c.OrderService.CancelEntireOrder(orderID, id); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Failure(err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, response.Success("Order cancelled successfully", nil))
}

// updating the order status
func (c *OrderController) UpdateOrderStatus(ctx *gin.Context) {
	orderID := ctx.Param("order_id")

	var dto dto.UpdateOrderStatusDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Failure("Invalid request body", err.Error()))
		return
	}

	// ADMIN PERMISSION CHECK
	role, _ := ctx.Get("UserRole")
	if role != "admin" {
		ctx.JSON(http.StatusForbidden, response.Failure("Forbidden: admin only", nil))
		return
	}

	if err := c.OrderService.UpdateOrderStatus(orderID, dto.Status); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Failure(err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, response.Success("Order status updated", nil))
}
