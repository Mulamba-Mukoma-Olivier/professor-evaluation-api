package eligibility

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

func (h *Handler) Check(c *gin.Context) {
	studentID, err := strconv.Atoi(c.Param("student_id"))

	if err != nil || studentID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid student ID",
		})
		return
	}

	eligibility, err := h.service.Check(studentID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, eligibility)
}