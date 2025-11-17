package routes

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.Engine) {
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
