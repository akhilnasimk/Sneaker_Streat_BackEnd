package routes

import (
	"github.com/akhilnasimk/SS_backend/internal/config"
	"github.com/akhilnasimk/SS_backend/internal/controllers"
	"github.com/akhilnasimk/SS_backend/internal/middlewares"
	"github.com/akhilnasimk/SS_backend/internal/repositories/sql"
	"github.com/akhilnasimk/SS_backend/internal/services"
	"github.com/gin-gonic/gin"
)

func RegisterWishlistRoute(rg gin.RouterGroup) {
	// Repositories
	wishrepo := sql.NewWishlistRepo(config.DB)
	productRepo := sql.NewProductsRepository(*config.DB)

	// Service
	wishService := services.NewWishlistService(wishrepo, productRepo)

	// Controller
	wishController := controllers.NewWishlistController(wishService)

	// Auth middlewares
	rg.Use(middlewares.AuthorizeMiddleware(), middlewares.CustomerAuth())
	{
		rg.GET("/", wishController.GetWishlist)
		rg.POST("/toggle/:product_id", wishController.ToggleWishlist)
		rg.DELETE("/:product_id", wishController.DeleteWishlistItem)
	}
}