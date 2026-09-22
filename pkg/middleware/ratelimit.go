package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     int
	window   time.Duration
}

type visitor struct {
	requests int
	lastSeen time.Time
}

func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		window:   window,
	}
	
	// Clean up old visitors periodically
	go rl.cleanup()
	
	return rl
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()
	
	for range ticker.C {
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > rl.window {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		
		rl.mu.Lock()
		v, exists := rl.visitors[ip]
		if !exists {
			rl.visitors[ip] = &visitor{requests: 1, lastSeen: time.Now()}
			rl.mu.Unlock()
			c.Next()
			return
		}
		
		if time.Since(v.lastSeen) > rl.window {
			v.requests = 1
			v.lastSeen = time.Now()
			rl.mu.Unlock()
			c.Next()
			return
		}
		
		v.requests++
		v.lastSeen = time.Now()
		rl.mu.Unlock()
		
		if v.requests > rl.rate {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}
