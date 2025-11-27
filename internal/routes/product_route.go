package routes

import (
	"github.com/akhilnasimk/SS_backend/internal/config"
	"github.com/akhilnasimk/SS_backend/internal/controllers"
	"github.com/akhilnasimk/SS_backend/internal/middlewares"
	"github.com/akhilnasimk/SS_backend/internal/repositories/sql"
	"github.com/akhilnasimk/SS_backend/internal/services"
	"github.com/gin-gonic/gin"
)

// RegisterProductRoutes sets up routes for products
func RegisterProductRoutes(rg *gin.RouterGroup) {

	// ---------------------
	// Repository Layer
	// ---------------------
	productRepo := sql.NewProductsRepository(*config.DB) // Product repository

	// ---------------------
	// Service Layer
	// ---------------------
	productService := services.NewProductsService(productRepo) // Product business logic

	// ---------------------
	// Controller Layer
	// ---------------------
	productController := controllers.NewProductController(productService)

	// ---------------------
	// Public Product Routes (Optional Auth)
	// ---------------------
	// OptionalAuth allows both logged-in users and guests to view products
	rg.GET("/", middlewares.OptionalAuth(), productController.GetAllProducts) // List all products
	rg.GET("/:id", productController.GetProductById)                          // Get product details by ID
	rg.GET("/categories", productController.GetAllCategory)                   // List all categories

	// ---------------------
	// Admin Product Routes (JWT + Admin Role)
	// ---------------------
	admin := rg.Group("/admin")
	admin.Use(middlewares.AuthorizeMiddleware(), middlewares.AdminAuth())
	{
		admin.POST("", productController.UploadProduct)                            // Add new product
		admin.PUT("/:id", productController.UpdateProduct)                          // Update product details
		admin.PATCH("/:id/toggle-availability", productController.ToggleProductAvailability) // Enable/disable product visibility
		admin.DELETE("/:id", productController.DeleteProduct)                        // Delete product
		admin.PATCH("/undelete/:id")                                                // Undelete soft-deleted product (if implemented)
	}
}
