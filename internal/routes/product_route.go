package routes

import (
	"github.com/akhilnasimk/SS_backend/internal/config"
	"github.com/akhilnasimk/SS_backend/internal/controllers"
	"github.com/akhilnasimk/SS_backend/internal/middlewares"
	"github.com/akhilnasimk/SS_backend/internal/repositories/sql"
	"github.com/akhilnasimk/SS_backend/internal/services"
	"github.com/gin-gonic/gin"
	// "myapp/controllers"
	// "myapp/middlewares"
)

// ProductRoutes sets up routes for products
func RegisterProductRoutes(rg *gin.RouterGroup) {
	//setting up the repo/service/controllers
	repo := sql.NewProductsRepository(*config.DB)
	Productservice := services.NewProductsService(repo)
	ProductController := controllers.NewProductController(Productservice)

	// Public product browsing
	rg.GET("/", ProductController.GetAllProducts)
	rg.GET("/:id", ProductController.GetProductById)

	//all the route that onyl accessable for the admin
	admin := rg.Group("/admin")
	admin.Use(middlewares.AuthorizeMiddleware(), middlewares.AdminAuth())
	{
		admin.POST("/", ProductController.UploadProduct)
	}

}
