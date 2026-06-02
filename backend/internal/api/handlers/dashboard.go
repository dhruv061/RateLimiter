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

// GetStats returns aggregated dashboard metrics scoped by global filters.
func (h *DashboardHandler) GetStats(c *gin.Context) {
	filter := parseGlobalFilter(c)
	stats, err := h.dashSvc.GetStats(filter)
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

// GetAttackStatus returns the current threat level scoped by global filters.
func (h *DashboardHandler) GetAttackStatus(c *gin.Context) {
	filter := parseGlobalFilter(c)
	status, err := h.dashSvc.GetAttackStatus(filter)
	if err != nil {
		response.InternalError(c, "Failed to fetch attack status")
		return
	}
	response.OK(c, status)
}
