package middleware

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// ValidateEmail validates email format
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidateMatricule validates matricule format (adjust pattern as needed)
func ValidateMatricule(matricule string) bool {
	// Example: alphanumeric, 6-20 characters
	matriculeRegex := regexp.MustCompile(`^[a-zA-Z0-9]{6,20}$`)
	return matriculeRegex.MatchString(matricule)
}

// SanitizeInput removes potentially dangerous characters
func SanitizeInput(input string) string {
	// Remove potential SQL injection patterns
	dangerousPatterns := []string{
		"'", ";", "--", "/*", "*/", "xp_", "sp_",
		"drop", "delete", "insert", "update", "exec",
	}

	sanitized := input
	for _, pattern := range dangerousPatterns {
		sanitized = strings.ReplaceAll(
			strings.ToLower(sanitized),
			strings.ToLower(pattern),
			"",
		)
	}

	return sanitized
}

// ValidateJSONContentType ensures request has JSON content type
func ValidateJSONContentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			contentType := c.GetHeader("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				c.JSON(http.StatusUnsupportedMediaType, gin.H{
					"error": "content-type must be application/json",
				})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

// ValidateRequestSize limits request body size
func ValidateRequestSize(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		c.Next()
	}
}
