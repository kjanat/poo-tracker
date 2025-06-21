package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter represents a rate limiter for API endpoints
type RateLimiter struct {
	mu      sync.RWMutex
	clients map[string]*ClientLimiter
	limit   int
	window  time.Duration
}

// ClientLimiter represents a rate limiter for a specific client
type ClientLimiter struct {
	mu       sync.Mutex
	requests []time.Time
	limit    int
	window   time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		clients: make(map[string]*ClientLimiter),
		limit:   limit,
		window:  window,
	}
}

// Allow checks if a request from the given client is allowed
func (rl *RateLimiter) Allow(clientID string) bool {
	rl.mu.RLock()
	client, exists := rl.clients[clientID]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		// Double-check locking pattern
		if client, exists = rl.clients[clientID]; !exists {
			client = &ClientLimiter{
				requests: make([]time.Time, 0),
				limit:    rl.limit,
				window:   rl.window,
			}
			rl.clients[clientID] = client
		}
		rl.mu.Unlock()
	}

	return client.Allow()
}

// Allow checks if a request is allowed for this client
func (cl *ClientLimiter) Allow() bool {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	now := time.Now()

	// Remove expired requests
	cutoff := now.Add(-cl.window)
	validRequests := cl.requests[:0]
	for _, req := range cl.requests {
		if req.After(cutoff) {
			validRequests = append(validRequests, req)
		}
	}
	cl.requests = validRequests

	// Check if we're under the limit
	if len(cl.requests) >= cl.limit {
		return false
	}

	// Add current request
	cl.requests = append(cl.requests, now)
	return true
}

// RateLimitMiddleware creates a Gin middleware for rate limiting
func RateLimitMiddleware(limiter *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use IP address as client ID, or user ID if authenticated
		clientID := c.ClientIP()
		if userID := c.GetString("userID"); userID != "" {
			clientID = userID
		}

		if !limiter.Allow(clientID) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// SecurityHeadersMiddleware adds security headers to responses
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self'")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		c.Next()
	}
}
