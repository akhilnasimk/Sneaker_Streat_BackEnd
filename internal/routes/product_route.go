package routes

import (
	// "github.com/akhilnasimk/SS_backend/internal/middlewares"
	"github.com/akhilnasimk/SS_backend/internal/config"
	"github.com/akhilnasimk/SS_backend/internal/controllers"
	"github.com/akhilnasimk/SS_backend/internal/repositories"
	"github.com/akhilnasimk/SS_backend/internal/services"
	"github.com/gin-gonic/gin"
	// "myapp/controllers"
	// "myapp/middlewares"
)

// ProductRoutes sets up routes for products
func RegisterProductRoutes(rg *gin.RouterGroup) {
	//setting up the repo/service/controllers
	repo := repositories.NewProductsRepository(*config.DB)
	Productservice := services.NewProductsService(repo)
	ProductController := controllers.NewProductController(Productservice)

	// Public product browsing
	rg.GET("/", ProductController.GetAllProducts)
	// rg.GET("/:id", GetProductByID)
	// // Optional filters
	// rg.GET("/category/:id", GetProductsByCategory)
	// rg.GET("/search/:query", SearchProducts)
}
