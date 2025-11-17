package middlewares

import (
	"net/http"
	"strings"

	"github.com/akhilnasimk/SS_backend/utils/jwt"
	"github.com/gin-gonic/gin"
)

func AuthorizeMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			ctx.Abort()
			return
		}

		// Remove "Bearer "
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		tokenString = strings.TrimSpace(tokenString)

		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
			ctx.Abort()
			return
		}

		// Validate token AND extract claims
		claims, err := jwt.ValidateAccessToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Message": "Access token validation failed",
				"Error":   err.Error(),
			})
			ctx.Abort()
			return
		}

		// Extract required data from claims
		userID := claims["UserId"]
		email := claims["UserEmail"]
		role := claims["UserRole"]


		// Store into Gin context
		ctx.Set("UserID", userID)
		ctx.Set("UserEmail", email)
		ctx.Set("UserRole", role)

		ctx.Next()
	}
}
