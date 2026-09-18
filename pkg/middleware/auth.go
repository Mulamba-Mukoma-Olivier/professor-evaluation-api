package middleware

import (
	"net/http"
	"strings"

	appjwt "github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func AuthRequired(jwtManager *appjwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Récupérer le header Authorization
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header is required",
			})
			c.Abort()
			return
		}

		// Vérifier le format : Bearer <token>
		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization format",
			})
			c.Abort()
			return
		}

		// Nettoyer le token
		token := strings.TrimSpace(parts[1])

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "token is required",
			})
			c.Abort()
			return
		}

		// Valider le JWT
		claims, err := jwtManager.ValidateToken(token)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			c.Abort()
			return
		}

		// Stocker les informations de l'utilisateur
		// dans le contexte Gin
		c.Set("user_id", claims.UserID)
		c.Set("matricule", claims.Matricule)
		c.Set("role", claims.Role)

		// Continuer vers le middleware/handler suivant
		c.Next()
	}
}