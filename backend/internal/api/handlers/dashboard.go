package handlers

import (
	"fail2ban-dashboard/internal/services"
	"fail2ban-dashboard/pkg/response"

	"github.com/gin-gonic/gin"
)

// DashboardHandler handles dashboard overview endpoints.
type DashboardHandler struct {
	dashSvc *services.DashboardService
}

// NewDashboardHandler creates a new DashboardHandler.
func NewDashboardHandler(dashSvc *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashSvc: dashSvc}
}

// GetStats returns aggregated dashboard metrics.
func (h *DashboardHandler) GetStats(c *gin.Context) {
	stats, err := h.dashSvc.GetStats()
	if err != nil {
		response.InternalError(c, "Failed to fetch dashboard stats")
		return
	}
	response.OK(c, stats)
}

// GetSystemStatus returns service health information.
func (h *DashboardHandler) GetSystemStatus(c *gin.Context) {
	status := h.dashSvc.GetSystemStatus()
	response.OK(c, status)
}

// GetAttackStatus returns the current threat level.
func (h *DashboardHandler) GetAttackStatus(c *gin.Context) {
	status, err := h.dashSvc.GetAttackStatus()
	if err != nil {
		response.InternalError(c, "Failed to fetch attack status")
		return
	}
	response.OK(c, status)
}
