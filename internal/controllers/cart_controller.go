package controllers

import (
	"fmt"
	"net/http"

	"github.com/akhilnasimk/SS_backend/internal/dto"
	"github.com/akhilnasimk/SS_backend/internal/enums"
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
	// Get userID from context
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

	// Call service
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

	// Return success
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
	orderitem, err := c.CartService.AddItemToCart(userID, productID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Failure("Failed to add product to cart", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.Success("Product added to cart successfully", orderitem))
}

func (c *CartController) UpdateCount(ctx *gin.Context) {
	//  Extract cart item ID
	cartItemId := ctx.Param("item_id")
	if cartItemId == "" {
		ctx.JSON(400, response.Failure("item_id is required", nil))
		return
	}

	// Extract op from query
	opStr := ctx.Query("op")
	if opStr == "" {
		ctx.JSON(400, response.Failure("op query parameter is required (inc/dec)", nil))
		return
	}

	//  Convert string → enum
	op := enums.CartOperation(opStr)

	//  Validate enum
	if !op.IsValid() {
		ctx.JSON(400, response.Failure("invalid op value (allowed: inc, dec)", nil))
		return
	}

	//  Call service
	if err := c.CartService.IncOrDecCartItem(cartItemId, opStr); err != nil {
		ctx.JSON(400, response.Failure("failed to update quantity", err.Error()))
		return
	}

	// 6️⃣ Success
	ctx.JSON(200, response.Success("quantity updated", nil))
}

func (c *CartController) DeleteCartItem(ctx *gin.Context) {
	id := ctx.Param("cartItemId")

	if id == "" {
		ctx.JSON(400, response.Failure("missing id as param", nil))
		return
	}

	if err := c.CartService.DeleteCartItem(id); err != nil {
		ctx.JSON(400, response.Failure("failed to delete the cart Items ", err.Error()))
		return
	}
	ctx.JSON(200, response.Success("successfully deleted product", nil))
}
