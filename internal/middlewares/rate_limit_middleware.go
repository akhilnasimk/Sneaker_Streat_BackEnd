package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// limiterStore keeps track of limiters per user/IP
type limiterStore struct {
	users map[string]*rate.Limiter
	mu    sync.Mutex
}

var store = &limiterStore{
	users: make(map[string]*rate.Limiter),
}

// GetLimiter returns the limiter for a given key (userID or IP)
func GetLimiter(key string) *rate.Limiter {
	store.mu.Lock()
	defer store.mu.Unlock()

	limiter, exists := store.users[key]
	if !exists {
		limiter = rate.NewLimiter(2, 10) // 5 req/sec, burst 10
		store.users[key] = limiter
	}
	return limiter
}

// RateLimitMiddleware limits requests per user/IP
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use IP address as key, can also use userID if authenticated
		clientIP := c.ClientIP()
		limiter := GetLimiter(clientIP)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
			})
			return
		}

		c.Next()
	}
}

// Optional: clean old limiters periodically
func CleanOldLimiters(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			store.mu.Lock()
			for _, limiter := range store.users {
				// remove limiters that have not been used recently
				// this requires adding a last-used timestamp in limiter struct (advanced)
				_ = limiter
			}
			store.mu.Unlock()
		}
	}()
}
