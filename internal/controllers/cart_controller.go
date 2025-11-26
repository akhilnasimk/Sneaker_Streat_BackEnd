package controllers

import (
	"fmt"
	"net/http"

	"github.com/akhilnasimk/SS_backend/internal/dto"
	"github.com/akhilnasimk/SS_backend/internal/helpers"
	"github.com/akhilnasimk/SS_backend/internal/services"
	"github.com/akhilnasimk/SS_backend/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CartController struct {
	CartService services.CartService
}

func NewCartController(cartService services.CartService) *CartController {
	return &CartController{
		CartService: cartService,
	}
}

// GetUserCart handles GET /cart
func (c *CartController) GetUserCart(ctx *gin.Context) {
	// 1️⃣ Get userID from context
	userIDValue, exists := ctx.Get("UserID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.Failure("User not authenticated", nil))
		return
	}

	userID := helpers.StringToUUID(userIDValue.(string))
	if userID == uuid.Nil {
		ctx.JSON(http.StatusBadRequest, response.Failure("Invalid user ID", nil))
		return
	}

	// 2️⃣ Call service
	cartResponse, err := c.CartService.GetUserCartItems(userID)
	if err != nil {
		// If record not found → return empty cart instead of failing
		if err.Error() == "record not found" {
			cartResponse = dto.CartResponse{
				Items: []dto.CartItemResponse{},
				Total: 0,
			}
		} else {
			ctx.JSON(http.StatusInternalServerError, response.Failure("Failed to fetch cart items", err.Error()))
			return
		}
	}

	// 3️⃣ Return success
	ctx.JSON(http.StatusOK, response.Success("Cart fetched successfully", cartResponse))
}

func (c *CartController) AddToCart(ctx *gin.Context) {
	//  Convert product ID from URL param
	productID := helpers.StringToUUID(ctx.Param("product_id"))
	if productID == uuid.Nil {
		ctx.JSON(http.StatusBadRequest, response.Failure("Invalid product ID", nil))
		return
	}

	//  Extract user ID from context (middleware sets it as string)
	userIDStr, exists := ctx.Get("UserID")
	fmt.Println("user Id is ", userIDStr)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.Failure("Unauthorized user", nil))
		return
	}

	// Convert userID (string) → UUID using your helper
	userID := helpers.StringToUUID(userIDStr.(string))
	if userID == uuid.Nil {
		ctx.JSON(http.StatusBadRequest, response.Failure("Invalid user ID", nil))
		return
	}

	// Service Layer call
	err := c.CartService.AddItemToCart(userID, productID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Failure("Failed to add product to cart", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.Success("Product added to cart successfully", nil))
}
