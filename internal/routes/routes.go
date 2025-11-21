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

	api := r.Group("/api/v1")
	// Auth routes: login, register, refresh
	auth := api.Group("/auth")
	AuthRoutes(auth)

	// Produts route all product and product by Id
	product := api.Group("/products")
	RegisterProductRoutes(product)

	users := api.Group("/users")
	RegisterUserRoutes(users)

}
