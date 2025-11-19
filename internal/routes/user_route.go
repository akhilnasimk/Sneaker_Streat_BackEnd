package routes

import (
	"github.com/akhilnasimk/SS_backend/internal/config"
	"github.com/akhilnasimk/SS_backend/internal/controllers"
	"github.com/akhilnasimk/SS_backend/internal/middlewares"
	"github.com/akhilnasimk/SS_backend/internal/repositories/sql"
	"github.com/akhilnasimk/SS_backend/internal/services"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup) {
	// Repository
	userRepo := sql.NewUserReposetory(*config.DB)
	// Service
	userService := services.NewUserService(userRepo)
	// Controller
	userController := controllers.NewUserController(userService)

	// ---------------------
	// JWT Protected Routes
	// ---------------------
	protected := rg.Group("/")
	protected.Use(middlewares.AuthorizeMiddleware()) // Validate JWT and set context
	{
		// ---------------------
		// Customer/User Routes
		// ---------------------
		customer := protected.Group("/")
		customer.Use(middlewares.CustomerAuth()) // Ensure role == "customer"
		{
			customer.GET("/me", userController.GetProfile)      // Get own profile
			customer.PATCH("/me", userController.UpdateProfile) // Update own profile
		}

		// ---------------------
		// Admin Routes
		// ---------------------
		admin := protected.Group("/admin")
		admin.Use(middlewares.AdminAuth()) // Ensure role == "admin"
		{
			admin.GET("/all-users", userController.GetAllUsers)       // Get all users
			admin.GET("/user/:id", userController.GetSingleUser)      // Get single user details
			admin.PATCH("/user/:id", userController.AdminUserUpdate) // Update user info or role
			admin.PUT("/users/:id/block", nil)                        // Block user account
		}
	}
}
