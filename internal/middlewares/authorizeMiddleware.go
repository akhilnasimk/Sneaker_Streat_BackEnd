package middlewares

import (
	"net/http"

	"github.com/akhilnasimk/SS_backend/utils/jwt"
	"github.com/gin-gonic/gin"
)

func AuthorizeMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := ctx.Cookie("access_token")
		if err != nil || tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Access token missing"})
			ctx.Abort()
			return
		}

		claims, err := jwt.ValidateAccessToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Access token validation failed",
				"error":   err.Error(),
			})
			ctx.Abort()
			return
		}

		ctx.Set("UserID", claims["UserId"])
		ctx.Set("UserEmail", claims["UserEmail"])
		ctx.Set("UserRole", claims["UserRole"])

		ctx.Next()
	}
}

