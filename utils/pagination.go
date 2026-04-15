package utils

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PaginationParams holds parsed pagination query params
type PaginationParams struct {
	Page     int
	PageSize int
	Offset   int
}

// GetPagination extracts and validates pagination params from query string
func GetPagination(c *gin.Context) PaginationParams {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Offset:   (page - 1) * pageSize,
	}
}

// BuildPaginatedResponse constructs a PaginatedResponse
func BuildPaginatedResponse(items interface{}, total int64, p PaginationParams) PaginatedResponse {
	totalPages := int(math.Ceil(float64(total) / float64(p.PageSize)))
	return PaginatedResponse{
		Items:      items,
		Total:      total,
		Page:       p.Page,
		PageSize:   p.PageSize,
		TotalPages: totalPages,
	}
}
