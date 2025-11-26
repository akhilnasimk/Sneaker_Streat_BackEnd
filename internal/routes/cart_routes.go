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

	//needed repos
	cartrepo := sql.NewcartRepository(*config.DB)
	productRepo := sql.NewProductsRepository(*config.DB)
	//needed services
	cartservice := services.NewCartService(cartrepo, productRepo)
	//needed controllers
	cartcontroller := controllers.NewCartController(cartservice)
	// All cart routes should require authentication
	rg.Use(middlewares.AuthorizeMiddleware())

	//  Get all cart items for current user
	rg.GET("/", cartcontroller.GetUserCart)
	// Add a product to cart
	rg.POST("/:product_id", cartcontroller.AddToCart)
	// // Update quantity of a cart item
	// rg.PATCH("/:id", controller.UpdateCartItem)
	// // Remove a product from cart
	// rg.DELETE("/:id", controller.RemoveFromCart)
	// // Optional: Clear entire cart for current user
	// rg.DELETE("/", controller.ClearCart)
}
