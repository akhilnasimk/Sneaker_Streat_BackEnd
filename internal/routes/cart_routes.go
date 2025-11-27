package routes

import (
	"github.com/akhilnasimk/SS_backend/internal/config"
	"github.com/akhilnasimk/SS_backend/internal/controllers"
	"github.com/akhilnasimk/SS_backend/internal/middlewares"
	"github.com/akhilnasimk/SS_backend/internal/repositories/sql"
	"github.com/akhilnasimk/SS_backend/internal/services"
	"github.com/gin-gonic/gin"
)

func RegisterCartRoutes(rg *gin.RouterGroup) {

	// ---------------------
	// Repository Layer
	// ---------------------
	cartRepo := sql.NewcartRepository(*config.DB)      // Cart repository
	productRepo := sql.NewProductsRepository(*config.DB) // Product repository (needed to fetch product details)

	// ---------------------
	// Service Layer
	// ---------------------
	cartService := services.NewCartService(cartRepo, productRepo) // Handles cart logic (add, update, delete)

	// ---------------------
	// Controller Layer
	// ---------------------
	cartController := controllers.NewCartController(cartService)

	// ---------------------
	// Cart Routes (JWT Protected)
	// ---------------------
	rg.Use(middlewares.AuthorizeMiddleware()) // Ensure user is authenticated
	{
		rg.GET("/", cartController.GetUserCart)                  // Get all cart items for current user
		rg.POST("/:product_id", cartController.AddToCart)        // Add a product to the cart
		rg.PATCH("/:item_id", cartController.UpdateCount)        // Increment/decrement quantity of a cart item
		rg.DELETE("/:cartItemId", cartController.DeleteCartItem) // Remove a product from the cart

		// Optional: Clear entire cart for current user
		// rg.DELETE("/", cartController.ClearCart)
	}
}
