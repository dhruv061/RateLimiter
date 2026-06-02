package handlers

import (
	"strconv"

	"fail2ban-dashboard/internal/services"
	"fail2ban-dashboard/pkg/response"

	"github.com/gin-gonic/gin"
)

// AnalyticsHandler handles traffic and security analytics endpoints.
type AnalyticsHandler struct {
	dashboardSvc *services.DashboardService
	banSvc       *services.BanService
}

// NewAnalyticsHandler creates a new AnalyticsHandler.
func NewAnalyticsHandler(dashboardSvc *services.DashboardService, banSvc *services.BanService) *AnalyticsHandler {
	return &AnalyticsHandler{dashboardSvc: dashboardSvc, banSvc: banSvc}
}

// GetTrafficTrends returns traffic trends scoped by global filters.
func (h *AnalyticsHandler) GetTrafficTrends(c *gin.Context) {
	filter := parseGlobalFilter(c)
	period := c.DefaultQuery("period", "hour")
	hours, _ := strconv.Atoi(c.DefaultQuery("hours", "24"))

	trends, err := h.dashboardSvc.GetTrafficTrends(filter, period, hours)
	if err != nil {
		response.InternalError(c, "Failed to fetch traffic trends: "+err.Error())
		return
	}

	response.OK(c, trends)
}

// GetCountryStats returns country-level security analytics scoped by global filters.
func (h *AnalyticsHandler) GetCountryStats(c *gin.Context) {
	filter := parseGlobalFilter(c)
	stats, err := h.dashboardSvc.GetCountryStats(filter)
	if err != nil {
		response.InternalError(c, "Failed to fetch country analytics: "+err.Error())
		return
	}

	response.OK(c, stats)
}

// GetTopOffenders returns top offending IPs scoped by global filters.
func (h *AnalyticsHandler) GetTopOffenders(c *gin.Context) {
	filter := parseGlobalFilter(c)
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offenders, err := h.banSvc.GetTopOffenders(filter, limit)
	if err != nil {
		response.InternalError(c, "Failed to fetch top offenders: "+err.Error())
		return
	}

	response.OK(c, offenders)
}
