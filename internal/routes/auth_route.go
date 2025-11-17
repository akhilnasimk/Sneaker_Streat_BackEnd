package routes

import (
	"github.com/akhilnasimk/SS_backend/internal/config"
	"github.com/akhilnasimk/SS_backend/internal/controllers"
	"github.com/akhilnasimk/SS_backend/internal/repositories"
	"github.com/akhilnasimk/SS_backend/internal/services"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(rg *gin.RouterGroup) {
	//All Repo that AuthController depends on
	userrepo := repositories.NewUserReposetory(*config.DB)   //creating new user repo
	tockenRepo := repositories.NewTokenRepository(config.DB) //creating a new Token repo
	otpRepo := repositories.NewOtpRepository(config.DB)      //creating a new Otp repo

	//ALL services that Auth controller depends on
	Authservice := services.NewAuthService(userrepo, tockenRepo) //auth service need userRepo and token repo
	EmailService := services.NewEmailService()                   //creating a Emailservice so put into Otp service
	OtpService := services.NewOtpService(otpRepo, EmailService)  //creating otp c with otp Repo and the Email service injected

	//the controller of the auth Route
	AuthController := controllers.NewAuthController(Authservice, OtpService) // finally creating with two services auth and otp

	auth := rg.Group("/")
	{
		auth.POST("/register", AuthController.Register)
		auth.POST("/login", AuthController.Login)
		auth.POST("/refresh", AuthController.RefreshToken)
		auth.POST("/forgotpassword", AuthController.ForgotPassword)
		auth.POST("/verify-otp", AuthController.VerifyOTP)
		auth.POST("/update-password", AuthController.UpdatePassword)
	}
}
