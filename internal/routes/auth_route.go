package routes

import (
	"github.com/akhilnasimk/SS_backend/internal/config"
	"github.com/akhilnasimk/SS_backend/internal/controllers"
	"github.com/akhilnasimk/SS_backend/internal/repositories/sql"
	"github.com/akhilnasimk/SS_backend/internal/services"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(rg *gin.RouterGroup) {

	// ---------------------
	// Repository Layer
	// ---------------------
	userRepo := sql.NewUserReposetory(*config.DB)  // User repository
	tokenRepo := sql.NewTokenRepository(config.DB) // Refresh token repository
	otpRepo := sql.NewOtpRepository(config.DB)     // OTP repository

	// ---------------------
	// Service Layer
	// ---------------------
	authService := services.NewAuthService(userRepo, tokenRepo) // Handles register/login/refresh
	emailService := services.NewEmailService()                  // Used by OTP service
	otpService := services.NewOtpService(otpRepo, emailService) // OTP generation/validation

	// ---------------------
	// Controller Layer
	// ---------------------
	authController := controllers.NewAuthController(authService, otpService)

	// ---------------------
	// Public Auth Routes
	// ---------------------
	auth := rg.Group("/")
	{
		auth.POST("/register", authController.Register)              // Register new user
		auth.POST("/login", authController.Login)                    // Login and get tokens
		auth.POST("/logout", authController.Logout)                  //Log out clear cookies and invalidate the refresh token in db
		auth.POST("/refresh", authController.RefreshToken)           // Refresh token
		auth.POST("/forgotpassword", authController.ForgotPassword)  // Send OTP to email
		auth.POST("/verify-otp", authController.VerifyOTP)           // Verify OTP
		auth.POST("/update-password", authController.UpdatePassword) // Update password after OTP verification
	}
}
