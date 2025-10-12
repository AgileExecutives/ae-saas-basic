package handlers

import (
	"net/http"
	"time"

	"github.com/ae-saas-basic/ae-saas-basic/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// Health performs a health check
// @Summary Health check
// @Description Check the health status of the API and database
// @Tags health
// @Produce json
// @Success 200 {object} models.HealthResponse
// @Failure 503 {object} models.ErrorResponse
// @Router /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	// Check database connection
	sqlDB, err := h.db.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, models.ErrorResponse{
			Error:   "Database connection error",
			Details: err.Error(),
		})
		return
	}

	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, models.ErrorResponse{
			Error:   "Database ping failed",
			Details: err.Error(),
		})
		return
	}

	response := models.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   "1.0.0", // TODO: Get from config
		Database:  "connected",
	}

	c.JSON(http.StatusOK, response)
}

// Ping performs a simple ping check
// @Summary Ping check
// @Description Simple ping endpoint
// @Tags health
// @Produce json
// @Success 200 {string} string "pong"
// @Router /ping [get]
func (h *HealthHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
