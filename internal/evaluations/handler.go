package evaluations

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EvaluationService interface {
	Create(studentID int, request CreateEvaluationRequest) (*Evaluation, error)
	GetAll() ([]Evaluation, error)
	GetByID(id int) (*Evaluation, error)
}

type Handler struct {
	service EvaluationService
}

func NewHandler(service EvaluationService) *Handler {
	return &Handler{
		service: service,
	}
}

// Create crée une évaluation pour l'étudiant authentifié.
func (h *Handler) Create(c *gin.Context) {
	var request CreateEvaluationRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	// Récupération de l'ID de l'étudiant depuis le JWT.
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	studentID, ok := userID.(int)
	if !ok || studentID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user ID",
		})
		return
	}

	evaluation, err := h.service.Create(studentID, request)

	if err != nil {
		if errors.Is(err, ErrEvaluationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "required resource not found",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, evaluation)
}

func (h *Handler) GetAll(c *gin.Context) {
	evaluations, err := h.service.GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"evaluations": evaluations,
	})
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid evaluation ID",
		})
		return
	}

	evaluation, err := h.service.GetByID(id)

	if err != nil {
		if errors.Is(err, ErrEvaluationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "evaluation not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, evaluation)
}
