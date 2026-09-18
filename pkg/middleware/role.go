package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Récupérer le rôle depuis le contexte
		roleValue, exists := c.Get("role")

		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "user role not found",
			})

			c.Abort()
			return
		}

		// Vérifier le type du rôle
		role, ok := roleValue.(string)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid user role",
			})

			c.Abort()
			return
		}

		// Normaliser le rôle
		role = strings.ToLower(strings.TrimSpace(role))

		// Vérifier les rôles autorisés
		for _, allowedRole := range allowedRoles {

			allowedRole = strings.ToLower(strings.TrimSpace(allowedRole))

			if role == allowedRole {
				c.Next()
				return
			}
		}

		// Rôle valide mais accès interdit
		c.JSON(http.StatusForbidden, gin.H{
			"error": "access denied",
		})

		c.Abort()
	}
}