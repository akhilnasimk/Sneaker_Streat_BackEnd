package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("UserRole")
		if !exists || role != "admin" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			ctx.Abort()
			return
		}

		// user is admin, pass through
		ctx.Next()
	}
}
