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

// NewRateLimiter crée un rate limiter.
// rate représente le nombre maximum de requêtes autorisées
// pendant une fenêtre de temps.
func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	if rate <= 0 {
		rate = 1
	}

	if window <= 0 {
		window = time.Minute
	}

	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		window:   window,
	}

	// Nettoyer périodiquement les visiteurs inactifs.
	go rl.cleanup()

	return rl
}

// cleanup supprime les visiteurs dont la fenêtre est expirée.
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()

		rl.mu.Lock()

		for ip, v := range rl.visitors {
			if now.Sub(v.lastSeen) > rl.window {
				delete(rl.visitors, ip)
			}
		}

		rl.mu.Unlock()
	}
}

// Middleware limite le nombre de requêtes par adresse IP.
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		rl.mu.Lock()

		v, exists := rl.visitors[ip]

		// Première requête de cette IP.
		if !exists {
			rl.visitors[ip] = &visitor{
				requests: 1,
				lastSeen: now,
			}

			rl.mu.Unlock()

			c.Next()
			return
		}

		// La fenêtre précédente est expirée.
		if now.Sub(v.lastSeen) > rl.window {
			v.requests = 1
			v.lastSeen = now

			rl.mu.Unlock()

			c.Next()
			return
		}

		// Nouvelle requête dans la fenêtre actuelle.
		v.requests++
		v.lastSeen = now

		requestCount := v.requests

		rl.mu.Unlock()

		// Limite dépassée.
		if requestCount > rl.rate {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
