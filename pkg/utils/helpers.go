package utils

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetPaginationParams extracts pagination parameters from query string
func GetPaginationParams(c *gin.Context) (page int, limit int) {
	page = 1
	limit = 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	return page, limit
}

// GetOffset calculates the database offset from page and limit
func GetOffset(page, limit int) int {
	return (page - 1) * limit
}

// CalculateTotalPages calculates total pages from total records and limit
func CalculateTotalPages(total int, limit int) int {
	if limit == 0 {
		return 0
	}
	return (total + limit - 1) / limit
}

// ValidateID validates and parses an ID parameter
func ValidateID(c *gin.Context, paramName string) (uint, error) {
	idStr := c.Param(paramName)
	if idStr == "" {
		return 0, fmt.Errorf("missing %s parameter", paramName)
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid %s parameter", paramName)
	}

	if id == 0 {
		return 0, fmt.Errorf("invalid %s parameter: must be greater than 0", paramName)
	}

	return uint(id), nil
}

// GetEnv gets environment variable with fallback
func GetEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
