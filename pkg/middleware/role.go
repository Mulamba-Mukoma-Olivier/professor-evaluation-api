package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// RequireRole vérifie que l'utilisateur possède l'un des rôles autorisés.
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Aucun rôle autorisé n'a été configuré.
		if len(allowedRoles) == 0 {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "no roles configured",
			})
			c.Abort()
			return
		}

		// Récupérer le rôle depuis le contexte Gin.
		roleValue, exists := c.Get("role")

		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "user role not found",
			})
			c.Abort()
			return
		}

		// DEBUG : afficher la valeur et son type.
		fmt.Printf(
			"DEBUG ROLE: %#v | TYPE: %T\n",
			roleValue,
			roleValue,
		)

		// Vérifier que le rôle est bien une chaîne.
		role, ok := roleValue.(string)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid user role",
			})
			c.Abort()
			return
		}

		// Normaliser le rôle.
		role = strings.ToLower(strings.TrimSpace(role))

		if role == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "user role is empty",
			})
			c.Abort()
			return
		}

		// DEBUG : afficher le rôle après normalisation
		// et les rôles autorisés.
		fmt.Printf(
			"DEBUG NORMALIZED ROLE: [%s]\n",
			role,
		)

		fmt.Printf(
			"DEBUG ALLOWED ROLES: %#v\n",
			allowedRoles,
		)

		// Vérifier si le rôle de l'utilisateur
		// correspond à l'un des rôles autorisés.
		for _, allowedRole := range allowedRoles {

			allowedRole = strings.ToLower(
				strings.TrimSpace(allowedRole),
			)

			if allowedRole == "" {
				continue
			}

			fmt.Printf(
				"DEBUG COMPARISON: [%s] == [%s]\n",
				role,
				allowedRole,
			)

			if role == allowedRole {
				fmt.Println("DEBUG ROLE AUTHORIZED")
				c.Next()
				return
			}
		}

		// L'utilisateur est authentifié mais
		// son rôle n'est pas autorisé.
		fmt.Println("DEBUG ROLE DENIED")

		c.JSON(http.StatusForbidden, gin.H{
			"error": "access denied",
		})
		c.Abort()
	}
}