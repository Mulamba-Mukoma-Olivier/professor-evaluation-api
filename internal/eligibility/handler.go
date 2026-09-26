package eligibility

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EligibilityService interface {
	Create(eligibility *Eligibility) (*Eligibility, error)
	Check(studentID int) (*Eligibility, error)
	Update(eligibility *Eligibility) (*Eligibility, error)
	Delete(studentID int) error
}

type Handler struct {
	service EligibilityService
}

func NewHandler(service EligibilityService) *Handler {
	return &Handler{
		service: service,
	}
}

// POST /eligibility
func (h *Handler) Create(c *gin.Context) {

	var request CreateEligibilityRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	if request.StudentID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid student ID",
		})
		return
	}

	eligibility := &Eligibility{
		StudentID:      request.StudentID,
		Enrollment:     request.Enrollment,
		AcademicFees:   request.AcademicFees,
		LaboratoryFees: request.LaboratoryFees,
		AccessFees:     request.AccessFees,
	}

	result, err := h.service.Create(eligibility)

	if err != nil {

		if errors.Is(err, ErrEligibilityNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GET /eligibility/:student_id
func (h *Handler) Check(c *gin.Context) {

	studentID, err := strconv.Atoi(
		c.Param("student_id"),
	)

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

// PUT /eligibility/:student_id
func (h *Handler) Update(c *gin.Context) {

	studentID, err := strconv.Atoi(
		c.Param("student_id"),
	)

	if err != nil || studentID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid student ID",
		})
		return
	}

	var request UpdateEligibilityRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	eligibility := &Eligibility{
		StudentID:      studentID,
		Enrollment:     request.Enrollment,
		AcademicFees:   request.AcademicFees,
		LaboratoryFees: request.LaboratoryFees,
		AccessFees:     request.AccessFees,
	}

	result, err := h.service.Update(eligibility)

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

	c.JSON(http.StatusOK, result)
}

// DELETE /eligibility/:student_id
func (h *Handler) Delete(c *gin.Context) {

	studentID, err := strconv.Atoi(
		c.Param("student_id"),
	)

	if err != nil || studentID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid student ID",
		})
		return
	}

	err = h.service.Delete(studentID)

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

	c.JSON(http.StatusOK, gin.H{
		"message": "eligibility deleted successfully",
	})
}