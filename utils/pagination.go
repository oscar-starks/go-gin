package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PaginationParams holds pagination parameters
type PaginationParams struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Offset   int `json:"offset"`
}

// PaginationResult holds pagination result metadata
type PaginationResult struct {
	Page         int   `json:"page"`
	PageSize     int   `json:"page_size"`
	TotalRecords int64 `json:"total_records"`
	TotalPages   int   `json:"total_pages"`
	HasNext      bool  `json:"has_next"`
	HasPrev      bool  `json:"has_prev"`
}

// GetPaginationParams extracts pagination parameters from Gin context
func GetPaginationParams(c *gin.Context) PaginationParams {
	page := 1
	pageSize := 10 // default page size

	// Get page from query parameter
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Get page_size from query parameter
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	offset := (page - 1) * pageSize

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Offset:   offset,
	}
}

// Paginate applies pagination to a GORM query and returns pagination metadata
func Paginate(db *gorm.DB, params PaginationParams) (*gorm.DB, PaginationResult) {
	var totalRecords int64

	// Count total records
	db.Count(&totalRecords)

	// Calculate pagination metadata
	totalPages := int((totalRecords + int64(params.PageSize) - 1) / int64(params.PageSize))
	hasNext := params.Page < totalPages
	hasPrev := params.Page > 1

	paginationResult := PaginationResult{
		Page:         params.Page,
		PageSize:     params.PageSize,
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
		HasNext:      hasNext,
		HasPrev:      hasPrev,
	}

	// Apply offset and limit to the query
	paginatedDB := db.Offset(params.Offset).Limit(params.PageSize)

	return paginatedDB, paginationResult
}
