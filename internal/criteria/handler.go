package criteria

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

// GET /criteria
func (h *Handler) GetAll(c *gin.Context) {
	criteria, err := h.service.GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrieve criteria",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"criteria": criteria,
	})
}

// GET /criteria/active
func (h *Handler) GetActive(c *gin.Context) {
	criteria, err := h.service.GetActive()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrieve active criteria",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"criteria": criteria,
	})
}

// GET /criteria/:id
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid criterion ID",
		})
		return
	}

	criterion, err := h.service.GetByID(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "criterion not found",
		})
		return
	}

	c.JSON(http.StatusOK, criterion)
}

// POST /criteria
func (h *Handler) Create(c *gin.Context) {
	var request CreateCriterionRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	criterion, err := h.service.Create(request)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, criterion)
}

// PUT /criteria/:id
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid criterion ID",
		})
		return
	}

	var request UpdateCriterionRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	criterion, err := h.service.Update(id, request)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, criterion)
}

// DELETE /criteria/:id
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid criterion ID",
		})
		return
	}

	err = h.service.Delete(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "criterion deleted successfully",
	})
}