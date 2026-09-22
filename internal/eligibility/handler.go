package eligibility

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EligibilityService interface {
	Check(studentID int) (*Eligibility, error)
}

type Handler struct {
	service EligibilityService
}

func NewHandler(service EligibilityService) *Handler {
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
		if errors.Is(err, ErrEligibilityNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "student eligibility not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, eligibility)
}
