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

// SecurityReport returns a compact security summary.
func (h *ReportsHandler) SecurityReport(c *gin.Context) {
	stats, err := h.dashboardSvc.GetStats()
	if err != nil {
		response.InternalError(c, "Failed to build security report")
		return
	}
	offenders, _ := h.banSvc.GetTopOffenders(10)
	logs, _ := h.auditSvc.GetRecentLogs(10)

	response.OK(c, map[string]interface{}{
		"generated_at": time.Now(),
		"stats":        stats,
		"top_attackers": offenders,
		"recent_audit": logs,
	})
}
