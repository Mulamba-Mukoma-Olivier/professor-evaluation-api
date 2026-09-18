package evaluations

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

func (h *Handler) Create(c *gin.Context) {
	studentID, err := strconv.Atoi(c.Param("student_id"))

	if err != nil || studentID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid student ID",
		})
		return
	}

	var request CreateEvaluationRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	evaluation, err := h.service.Create(studentID, request)

	if err != nil {
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
			"error": err.Error(),
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
		c.JSON(http.StatusNotFound, gin.H{
			"error": "evaluation not found",
		})
		return
	}

	c.JSON(http.StatusOK, evaluation)
}