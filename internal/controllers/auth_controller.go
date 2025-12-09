package controllers

import (
	"log"
	"net/http"

	"github.com/akhilnasimk/SS_backend/internal/dto"
	"github.com/akhilnasimk/SS_backend/internal/helpers"
	"github.com/akhilnasimk/SS_backend/internal/services"
	"github.com/akhilnasimk/SS_backend/utils/jwt"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService services.AuthService
	OtpService  services.OtpService
}

func NewAuthController(service services.AuthService, OtpS services.OtpService) *AuthController {
	return &AuthController{
		authService: service,
		OtpService:  OtpS,
	}
}

func (r *AuthController) Register(ctx *gin.Context) {

	var user dto.RegisterRequest

	// Bind JSON input
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body: " + err.Error(),
		})
		return
	}
	log.Println(user)
	if err := helpers.ValidateUserInput(user.UserName, user.Email, user.Password, *user.Phone); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "user validation failed ",
			"error":   err.Error(),
		})
		return
	}

	// Call service
	err := r.authService.Register(user)
	if err != nil {
		// Check for user already exists error
		if err.Error() == "user with email "+user.Email+" already exists" {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Handle bcrypt or DB errors
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Success
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "user has been registered successfully",
		"user": gin.H{
			"username": user.UserName,
			"email":    user.Email,
		},
	})
}

func (r *AuthController) Login(ctx *gin.Context) {
	var loginReq dto.LoginReq

	if err := ctx.ShouldBindJSON(&loginReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// 2. Authenticate
	user, err := r.authService.Login(loginReq)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 3. Generate Access Token
	accessToken, err := jwt.GenerateAccess(user.ID, user.UserName, user.Email, *user.UserRole)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	// 4. Generate Refresh Token
	refreshToken, err := r.authService.GenerateAndStoreRefreshToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	// 5. Set HTTP-Only Cookies
	// Short-lived Access Token cookie (optional, can also stay in memory)
	ctx.SetCookie(
		"access_token",
		accessToken,
		15*60, // 15 minutes
		"/",
		"",
		false, // secure=false in dev, true in production
		true,  // httpOnly
	)

	// Long-lived Refresh Token cookie
	ctx.SetCookie(
		"refresh_token",
		refreshToken,
		7*24*60*60, // 7 days
		"/",
		"",
		false,
		true,
	)

	// Response
	// Send user info in response
	userInfo := gin.H{
		"id":       user.ID,
		"username": user.UserName,
		"email":    user.Email,
		"role":     user.UserRole,
		"image":    user.Image,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"user":    userInfo,
	})
}

func (r *AuthController) Logout(ctx *gin.Context) {
	// 1. Get refresh token from cookie
	refreshToken, err := ctx.Cookie("refresh_token")
	if err == nil && refreshToken != "" {
		// Invalidate refresh token in DB
		_ = r.authService.InvalidateRefreshToken(refreshToken)
	}

	// 2. Delete Cookies (overwrite with empty + expired)
	ctx.SetCookie("access_token", "", -1, "/", "", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "", false, true)

	// 3. Response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "logged out successfully",
	})
}

func (r *AuthController) RefreshToken(ctx *gin.Context) {
	refresh, err := ctx.Cookie("refresh_token")
	if err != nil || refresh == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Refresh token missing", "error": err.Error()})
		return
	}

	// Generate new access and refresh tokens
	accessToken, newRefreshToken, err := r.authService.RefreshTokens(refresh)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Failed to refresh tokens", "error": err.Error()})
		return
	}

	// Set new cookies
	ctx.SetCookie(
		"access_token",
		accessToken,
		15*60, // 15 minutes
		"/",
		"",
		false,
		true,
	)

	ctx.SetCookie(
		"refresh_token",
		newRefreshToken,
		7*24*60*60,
		"/",
		"",
		false,
		true,
	)

	ctx.JSON(http.StatusOK, gin.H{"message": "Tokens refreshed"})
}

// forgot passwrd with the otp sending
func (r *AuthController) ForgotPassword(ctx *gin.Context) {
	var req dto.ForgotPasswordReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad request",
			"error":   err.Error(),
		})
		return
	}

	//checkign for user
	user, err := r.authService.ForgotPassword(req.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "User does not exist or is blocked",
			"error":   err.Error(),
		})
		return
	}

	// otp with our otp service
	if err := r.OtpService.SendOTP(user.ID, user.Email, "forgot_password"); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to send OTP",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "OTP sent to your email",
		"email":   user.Email,
	})
}

func (r *AuthController) VerifyOTP(ctx *gin.Context) {
	var otpCheckReq dto.VerifyOTPReq

	if err := ctx.ShouldBindBodyWithJSON(&otpCheckReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad request",
			"error":   err.Error(),
		})
		return
	}

	verified, err := r.OtpService.VerifyOTP(otpCheckReq.Otp, otpCheckReq.Email, "forgot_password")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "something went wrong",
			"error":   err.Error(),
		})
		return
	}

	if !verified {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "otp verification failed ",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"message": "Otp is verified",
		"email":   otpCheckReq.Email,
	})
}

// update password controller
func (r *AuthController) UpdatePassword(ctx *gin.Context) {
	var Request dto.ChangePasswordReq

	if err := ctx.ShouldBindBodyWithJSON(&Request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "binding failed ",
			"error":   err.Error(),
		})
		return
	}
	//running the methode from the service
	if err := r.authService.UpdatePassword(Request.Email, Request.NewPassword); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to update the password ",
			"err":     err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Pass word has been changed ",
		"email":   Request.Email,
	})
}
