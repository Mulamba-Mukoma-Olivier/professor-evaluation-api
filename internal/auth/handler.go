package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	appjwt "github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/pkg/jwt"
)

type Handler struct {
	service *Service
	jwt     *appjwt.Manager
}

func NewHandler(service *Service, jwt *appjwt.Manager) *Handler {
	return &Handler{
		service: service,
		jwt:     jwt,
	}
}

// Login authentifie un utilisateur et retourne un JWT.
func (h *Handler) Login(c *gin.Context) {
	var request LoginRequest

	// Validation du JSON
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	// Authentification
	user, err := h.service.Login(request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid credentials",
		})
		return
	}

	// Génération du JWT
	token, err := h.jwt.GenerateToken(
		user.ID,
		user.Matricule,
		user.Role,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate token",
		})
		return
	}

	// Réponse
	resp := LoginResponse{
		AccessToken: token,
		TokenType:   "bearer",
		User:        *user,
	}

	c.JSON(http.StatusOK, resp)
}

// Register crée un nouvel utilisateur.
func (h *Handler) Register(c *gin.Context) {
	var request RegisterRequest

	// Validation du JSON
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	// Création de l'utilisateur
	user, err := h.service.Register(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Réponse
	c.JSON(http.StatusCreated, gin.H{
		"message": "user created",
		"user":    user,
	})
}
