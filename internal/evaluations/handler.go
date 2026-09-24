package evaluations

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EvaluationService interface {
	Create(
		studentID int,
		request CreateEvaluationRequest,
	) (*Evaluation, error)

	GetAll() ([]Evaluation, error)

	GetByID(id int) (*Evaluation, error)

	Delete(id int) error
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

	// --------------------------------------------------
	// Validation du JSON
	// --------------------------------------------------
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	// --------------------------------------------------
	// Récupération de l'utilisateur depuis le JWT
	// --------------------------------------------------
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

	// --------------------------------------------------
	// Création de l'évaluation
	// --------------------------------------------------
	evaluation, err := h.service.Create(studentID, request)

	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidStudentID),
			errors.Is(err, ErrInvalidProfessorID),
			errors.Is(err, ErrInvalidCourseID),
			errors.Is(err, ErrAcademicYearRequired),
			errors.Is(err, ErrPeriodRequired),
			errors.Is(err, ErrAnswersRequired),
			errors.Is(err, ErrInvalidCriterionID),
			errors.Is(err, ErrInvalidScore),
			errors.Is(err, ErrDuplicateCriterion),
			errors.Is(err, ErrDuplicateEvaluation):

			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			return
		}
	}

	c.JSON(http.StatusCreated, evaluation)
}

// GetAll retourne toutes les évaluations.
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

// GetByID retourne une évaluation par son ID.
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
		if errors.Is(err, ErrInvalidEvaluationID) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

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

// Delete supprime une évaluation par son ID.
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid evaluation ID",
		})
		return
	}

	err = h.service.Delete(id)

	if err != nil {
		if errors.Is(err, ErrInvalidEvaluationID) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

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

	c.JSON(http.StatusOK, gin.H{
		"message": "evaluation deleted successfully",
	})
}
