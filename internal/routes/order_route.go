package routes

import (
	"github.com/akhilnasimk/SS_backend/internal/config"
	"github.com/akhilnasimk/SS_backend/internal/controllers"
	"github.com/akhilnasimk/SS_backend/internal/middlewares"
	"github.com/akhilnasimk/SS_backend/internal/repositories/sql"
	"github.com/akhilnasimk/SS_backend/internal/services"
	"github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(rg *gin.RouterGroup) {
	//repositories
	OrderRepo := sql.NewOrderRepository(*config.DB)
	Cartrepo := sql.NewcartRepository(*config.DB)

	//services
	OrderService := services.NewOrderService(OrderRepo, Cartrepo)

	//controller
	OrderController := controllers.NewOrderController(OrderService)

	rg.Use(middlewares.AuthorizeMiddleware())
	{
		// Get all orders for logged-in user
		rg.GET("/", OrderController.GetAllOrders)

		// Create orders
		rg.POST("/single/:product_id", OrderController.AddSingleItemOrder)
		rg.POST("/cart/:cart_id", OrderController.AddCartOrder)

		// Cancel operations
		rg.DELETE("/items/:item_id/cancel", OrderController.CancelOrderItem)
		rg.DELETE("/:order_id/cancel", OrderController.CancelOrder)
	}
	admin := rg.Group("/admin")
	admin.Use(middlewares.AdminAuth())
	{
		admin.PATCH("/:order_id/status", OrderController.UpdateOrderStatus)
	}

}
