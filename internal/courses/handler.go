package courses

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CourseService interface {
	GetAll() ([]Course, error)
	GetByID(id int) (*Course, error)
	Create(request CreateCourseRequest) (*Course, error)
	Update(id int, request UpdateCourseRequest) (*Course, error)
	Delete(id int) error
}

type Handler struct {
	service CourseService
}

func NewHandler(service CourseService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetAll(c *gin.Context) {
	courses, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"courses": courses,
	})
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid course ID",
		})
		return
	}

	course, err := h.service.GetByID(id)
	if err != nil {
		if errors.Is(err, ErrCourseNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "course not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, course)
}

func (h *Handler) Create(c *gin.Context) {
	var request CreateCourseRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	course, err := h.service.Create(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, course)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid course ID",
		})
		return
	}

	var request UpdateCourseRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	course, err := h.service.Update(id, request)
	if err != nil {
		if errors.Is(err, ErrCourseNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "course not found",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, course)
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid course ID",
		})
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		if errors.Is(err, ErrCourseNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "course not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "course deleted successfully",
	})
}
