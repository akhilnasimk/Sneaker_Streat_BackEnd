package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CustomerAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := ctx.Get("UserID")
		if !exists || userID == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			ctx.Abort()
			return
		}

		role, exists := ctx.Get("UserRole")
		if !exists || role != "customer" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Customer access required"})
			ctx.Abort()
			return
		}

		// Authenticated as a customer
		ctx.Next()
	}
}
