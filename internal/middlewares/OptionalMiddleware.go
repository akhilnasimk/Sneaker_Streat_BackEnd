package middlewares

import (
	"github.com/akhilnasimk/SS_backend/utils/jwt"
	"github.com/gin-gonic/gin"
)

func OptionalAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// fmt.Println("no token valid")
		tokenString, err := ctx.Cookie("access_token")
		if err != nil || tokenString == "" {
			// No cookie → treat as guest
			ctx.Next()
			return
		}

		claims, err := jwt.ValidateAccessToken(tokenString)
		if err == nil {
			// Token valid → set user values
			ctx.Set("UserID", claims["UserId"])
			ctx.Set("UserEmail", claims["UserEmail"])
			ctx.Set("UserRole", claims["UserRole"])
			// fmt.Println("role has beeen set ")
		}

		// Even if token invalid → do NOT block
		// OptionalAuth should NEVER abort
		ctx.Next()
	}
}
