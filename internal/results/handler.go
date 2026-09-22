package results

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ResultsService interface {
	GetProfessorResult(
		professorID int,
		courseID int,
		academicYear string,
		period string,
	) (*ProfessorResult, error)
}

type Handler struct {
	service ResultsService
}

func NewHandler(service ResultsService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetProfessorResult(c *gin.Context) {

	// Récupération de l'ID du professeur
	professorID, err := strconv.Atoi(c.Param("professor_id"))

	if err != nil || professorID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid professor ID",
		})
		return
	}

	// Récupération de l'ID du cours
	courseID, err := strconv.Atoi(c.Query("course_id"))

	if err != nil || courseID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid course ID",
		})
		return
	}

	// Année académique
	academicYear := strings.TrimSpace(c.Query("academic_year"))

	if academicYear == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "academic_year is required",
		})
		return
	}

	// Période
	period := strings.TrimSpace(c.Query("period"))

	if period == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "period is required",
		})
		return
	}

	// Appel du service
	result, err := h.service.GetProfessorResult(
		professorID,
		courseID,
		academicYear,
		period,
	)

	if err != nil {

		// Aucune évaluation trouvée
		if errors.Is(err, ErrNoEvaluations) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Les évaluations existent mais ne contiennent aucune réponse
		if errors.Is(err, ErrNoAnswers) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Erreur interne
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
