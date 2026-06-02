package handlers

import (
	"time"

	"fail2ban-dashboard/internal/services"
	"fail2ban-dashboard/pkg/response"

	"github.com/gin-gonic/gin"
)

// ReportsHandler handles report endpoints.
type ReportsHandler struct {
	dashboardSvc *services.DashboardService
	banSvc       *services.BanService
	auditSvc     *services.AuditService
}

// NewReportsHandler creates a new ReportsHandler.
func NewReportsHandler(dashboardSvc *services.DashboardService, banSvc *services.BanService, auditSvc *services.AuditService) *ReportsHandler {
	return &ReportsHandler{dashboardSvc: dashboardSvc, banSvc: banSvc, auditSvc: auditSvc}
}

// SecurityReport returns a compact security summary scoped by global filters.
func (h *ReportsHandler) SecurityReport(c *gin.Context) {
	filter := parseGlobalFilter(c)
	stats, err := h.dashboardSvc.GetStats(filter)
	if err != nil {
		response.InternalError(c, "Failed to build security report: "+err.Error())
		return
	}
	offenders, _ := h.banSvc.GetTopOffenders(filter, 10)
	logs, _ := h.auditSvc.GetRecentLogs(filter.DomainID, 10)

	response.OK(c, map[string]interface{}{
		"generated_at": time.Now(),
		"stats":        stats,
		"top_attackers": offenders,
		"recent_audit":  logs,
	})
}
