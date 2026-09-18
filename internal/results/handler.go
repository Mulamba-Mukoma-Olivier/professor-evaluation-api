package results

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetProfessorResult(c *gin.Context) {

	// Récupération de l'ID du professeur
	professorID, err := strconv.Atoi(
		c.Param("professor_id"),
	)

	if err != nil || professorID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid professor ID",
		})
		return
	}

	// Récupération de l'ID du cours
	courseID, err := strconv.Atoi(
		c.Query("course_id"),
	)

	if err != nil || courseID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid course ID",
		})
		return
	}

	// Année académique
	academicYear := c.Query("academic_year")

	if academicYear == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "academic_year is required",
		})
		return
	}

	// Période
	period := c.Query("period")

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
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}