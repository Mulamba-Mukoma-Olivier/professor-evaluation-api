package professors

import (
	"net/http"
	"strconv"

	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/pkg/pagination"
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

func (h *Handler) GetAll(c *gin.Context) {
	page, pageSize := pagination.GetPagination(c)
	
	professors, total, err := h.service.GetAllPaginated(page, pageSize)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := pagination.BuildResponse(professors, page, pageSize, total)
	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetActive(c *gin.Context) {
	professors, err := h.service.GetActive()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"professors": professors,
	})
}

func (h *Handler) GetByStatus(c *gin.Context) {
	status := c.Query("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "status parameter is required",
		})
		return
	}

	professors, err := h.service.GetByStatus(status)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"professors": professors,
	})
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid professor ID",
		})
		return
	}

	professor, err := h.service.GetByID(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "professor not found",
		})
		return
	}

	c.JSON(http.StatusOK, professor)
}

func (h *Handler) Create(c *gin.Context) {
	var request CreateProfessorRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	professor, err := h.service.Create(request)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, professor)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid professor ID",
		})
		return
	}

	var request UpdateProfessorRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	professor, err := h.service.Update(id, request)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, professor)
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid professor ID",
		})
		return
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "professor deleted successfully",
	})
}