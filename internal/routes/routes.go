package routes

import (
	_ "github.com/akhilnasimk/SS_backend/docs"
	"github.com/akhilnasimk/SS_backend/internal/middlewares"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(r *gin.Engine) {
	// base middlewares
	r.Use(middlewares.CORSMiddleware())
	r.Use(middlewares.RateLimitMiddleware())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//for testing purpose only 
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	api := r.Group("/api/v1")
	// Auth routes: login, register, refresh
	auth := api.Group("/auth")
	AuthRoutes(auth)

	// Produts route all related to products
	product := api.Group("/products")
	RegisterProductRoutes(product)

	//routes that is related to users
	users := api.Group("/users")
	RegisterUserRoutes(users)

	//route that are related to the cart service
	cart := api.Group("/cart")
	RegisterCartRoutes(cart)

	//route that are related to wishlist  service
	wishlist := api.Group("/wishlist")
	RegisterWishlistRoute(*wishlist)

	//route that are realted to the orders
	Order := api.Group("/order")
	RegisterOrderRoutes(Order)

}
