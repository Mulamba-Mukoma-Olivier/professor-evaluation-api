package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewRateLimiter(t *testing.T) {
	rl := NewRateLimiter(10, time.Minute)
	
	assert.NotNil(t, rl)
	assert.Equal(t, 10, rl.rate)
	assert.Equal(t, time.Minute, rl.window)
}

func TestRateLimiter_Middleware_AllowRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rl := NewRateLimiter(5, time.Minute)
	
	router := gin.New()
	router.Use(rl.Middleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	
	// Make 5 requests (should all succeed)
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "127.0.0.1"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
	}
}

func TestRateLimiter_Middleware_BlockExcessRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rl := NewRateLimiter(2, time.Minute)
	
	router := gin.New()
	router.Use(rl.Middleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	
	// Make 2 requests (should succeed)
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "127.0.0.1"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
	}
	
	// Make 3rd request (should be blocked)
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1"
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
}

func TestRateLimiter_Middleware_DifferentIPs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rl := NewRateLimiter(2, time.Minute)
	
	router := gin.New()
	router.Use(rl.Middleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	
	// First IP - should succeed
	req1 := httptest.NewRequest("GET", "/test", nil)
	req1.RemoteAddr = "127.0.0.1:12345"
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)
	
	// Second IP - should also succeed
	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.RemoteAddr = "127.0.0.2:12345"
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}
