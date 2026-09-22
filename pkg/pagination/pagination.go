package pagination

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 100
)

type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// GetPagination récupère et valide les paramètres de pagination.
func GetPagination(c *gin.Context) (int, int) {
	page, err := strconv.Atoi(
		c.DefaultQuery("page", strconv.Itoa(DefaultPage)),
	)
	if err != nil || page < 1 {
		page = DefaultPage
	}

	pageSize, err := strconv.Atoi(
		c.DefaultQuery("page_size", strconv.Itoa(DefaultPageSize)),
	)
	if err != nil || pageSize < 1 || pageSize > MaxPageSize {
		pageSize = DefaultPageSize
	}

	return page, pageSize
}

// GetOffset calcule l'offset SQL correspondant à la page.
func GetOffset(page, pageSize int) int {
	if page < 1 {
		page = DefaultPage
	}

	if pageSize < 1 {
		pageSize = DefaultPageSize
	}

	return (page - 1) * pageSize
}

// GetTotalPages calcule le nombre total de pages.
func GetTotalPages(total, pageSize int) int {
	if total <= 0 || pageSize <= 0 {
		return 0
	}

	return int(math.Ceil(float64(total) / float64(pageSize)))
}

// BuildResponse construit une réponse paginée.
func BuildResponse(
	data interface{},
	page int,
	pageSize int,
	total int,
) PaginatedResponse {
	return PaginatedResponse{
		Data: data,
		Pagination: Pagination{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: GetTotalPages(total, pageSize),
		},
	}
}
