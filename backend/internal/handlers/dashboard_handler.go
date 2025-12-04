package handlers

import (
	"net/http"
	"github.com/Mahfuz2811/medecole/backend/internal/models"
	"github.com/Mahfuz2811/medecole/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// DashboardHandler handles dashboard-related endpoints
type DashboardHandler struct {
	dashboardService service.DashboardService
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler(dashboardService service.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

func (h *DashboardHandler) GetDashboardSummary(c *gin.Context) {
	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "Unauthorized",
			Message: "Invalid user ID",
		})
		return
	}

	// Get fresh dashboard data directly from database
	response, err := h.dashboardService.GetDashboardSummary(uid)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to retrieve dashboard data: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Dashboard data retrieved successfully",
		"data":    response,
	})
}

func (h *DashboardHandler) GetDashboardEnrollments(c *gin.Context) {
	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Invalid user ID format",
		})
		return
	}

	// Get dashboard enrollments from service
	enrollments, err := h.dashboardService.GetDashboardEnrollments(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to retrieve enrollments: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, enrollments)
}
