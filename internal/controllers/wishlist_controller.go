package controllers

import (
	"net/http"

	"github.com/akhilnasimk/SS_backend/internal/services"
	"github.com/akhilnasimk/SS_backend/utils/response"
	"github.com/gin-gonic/gin"
)

type WishlistController struct {
	wishlistService services.WishlistService
}

func NewWishlistController(wishlistService services.WishlistService) *WishlistController {
	return &WishlistController{
		wishlistService: wishlistService,
	}
}

func (c *WishlistController) GetWishlist(ctx *gin.Context) {
	userID, exists := ctx.Get("UserID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.Failure("user not authorized", nil))
		return
	}

	items, err := c.wishlistService.GetAllWishlistItems(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Failure("failed to fetch wishlist items", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.Success("wishlist fetched successfully", items))
}

func (c *WishlistController) ToggleWishlist(ctx *gin.Context) {
	userID, exists := ctx.Get("UserID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.Failure("user not authorized", nil))
		return
	}

	productID := ctx.Param("product_id")
	if productID == "" {
		ctx.JSON(http.StatusBadRequest, response.Failure("product_id is required", nil))
		return
	}

	action, item, err := c.wishlistService.ToggleWishlist(userID.(string), productID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Failure(err.Error(), nil))
		return
	}

	// Response DTO
	resp := gin.H{
		"action": action, // "added" or "removed"
		"item":   item,
	}

	ctx.JSON(http.StatusOK, response.Success("wishlist updated", resp))
}

func (c *WishlistController) DeleteWishlistItem(ctx *gin.Context) {
	userID, exists := ctx.Get("UserID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.Failure("user not authorized", nil))
		return
	}

	productID := ctx.Param("product_id")
	if productID == "" {
		ctx.JSON(http.StatusBadRequest, response.Failure("product_id is required", nil))
		return
	}

	err := c.wishlistService.DeleteWishlistItem(userID.(string), productID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Failure(err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, response.Success("wishlist item deleted successfully", nil))
}

func (c *WishlistController) CheckProduct(ctx *gin.Context) {
	userId, exist := ctx.Get("UserID")
	if !exist {
		ctx.JSON(400, response.Failure("Un Autherized ", nil))
		return
	}

	PId := ctx.Param("product-id")
	if PId == "" {
		ctx.JSON(400, response.Failure("Need to pass the product ID  ", nil))
		return
	}

	resp, err := c.wishlistService.CheckWishlistStatus(userId.(string), PId)

	if err != nil {
		ctx.JSON(400, response.Failure("Something went wrong in the check ", err))
		return
	}

	ctx.JSON(200, response.Success("Check sucess", resp))
}
