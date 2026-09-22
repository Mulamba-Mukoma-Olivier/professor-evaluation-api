package middleware

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	emailRegex     = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$`)
	matriculeRegex = regexp.MustCompile(`^[a-zA-Z0-9]{6,20}$`)
)

// ValidateEmail vérifie le format général d'une adresse email.
func ValidateEmail(email string) bool {
	email = strings.TrimSpace(email)

	if email == "" || len(email) > 254 {
		return false
	}

	return emailRegex.MatchString(email)
}

// ValidateMatricule vérifie le format d'un matricule.
// Format accepté : 6 à 20 caractères alphanumériques.
func ValidateMatricule(matricule string) bool {
	matricule = strings.TrimSpace(matricule)

	return matriculeRegex.MatchString(matricule)
}

// SanitizeInput nettoie uniquement les espaces inutiles.
// La protection contre les injections SQL doit être assurée
// par l'utilisation de requêtes paramétrées / GORM.
func SanitizeInput(input string) string {
	return strings.TrimSpace(input)
}

// ValidateJSONContentType vérifie que les requêtes
// POST, PUT et PATCH utilisent application/json.
func ValidateJSONContentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch:
			contentType := c.GetHeader("Content-Type")

			if contentType == "" ||
				!strings.HasPrefix(
					strings.ToLower(contentType),
					"application/json",
				) {
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

// ValidateRequestSize limite la taille du corps de la requête.
func ValidateRequestSize(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if maxSize <= 0 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "invalid maximum request size",
			})
			c.Abort()
			return
		}

		c.Request.Body = http.MaxBytesReader(
			c.Writer,
			c.Request.Body,
			maxSize,
		)

		c.Next()
	}
}
