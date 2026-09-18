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

func (h *Handler) Login(c *gin.Context) {

	var request LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	user, err := h.service.Login(request)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid credentials",
		})
		return
	}

	// generate token
	token, err := h.jwt.GenerateToken(user.ID, user.Matricule, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	resp := LoginResponse{
		AccessToken: token,
		TokenType:   "bearer",
		User:        *user,
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) Register(c *gin.Context) {
	var request RegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := h.service.Register(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created",
		"user":    user,
	})
}