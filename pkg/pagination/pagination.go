package pagination

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

func GetPagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return page, pageSize
}

func GetOffset(page, pageSize int) int {
	return (page - 1) * pageSize
}

func GetTotalPages(total, pageSize int) int {
	if pageSize == 0 {
		return 0
	}
	return int(math.Ceil(float64(total) / float64(pageSize)))
}

func BuildResponse(data interface{}, page, pageSize, total int) PaginatedResponse {
	_ = GetTotalPages(total, pageSize) // Calculate but not used for now

	return PaginatedResponse{
		Data: data,
		Pagination: Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	}
}
