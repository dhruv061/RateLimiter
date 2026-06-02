package handlers

import (
	"strconv"

	"fail2ban-dashboard/internal/services"
	"fail2ban-dashboard/pkg/response"

	"github.com/gin-gonic/gin"
)

// AuditHandler handles audit log endpoints.
type AuditHandler struct {
	auditSvc *services.AuditService
}

// NewAuditHandler creates a new AuditHandler.
func NewAuditHandler(auditSvc *services.AuditService) *AuditHandler {
	return &AuditHandler{auditSvc: auditSvc}
}

// GetLogs returns paginated audit logs scoped by global filters.
func (h *AuditHandler) GetLogs(c *gin.Context) {
	filter := parseGlobalFilter(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	search := c.Query("search")
	action := c.Query("action")

	result, err := h.auditSvc.GetLogs(filter, page, perPage, search, action)
	if err != nil {
		response.InternalError(c, "Failed to fetch audit logs: "+err.Error())
		return
	}
	response.OK(c, result)
}

// Export returns audit logs for export scoped by global filters.
func (h *AuditHandler) Export(c *gin.Context) {
	filter := parseGlobalFilter(c)
	result, err := h.auditSvc.GetLogs(filter, 1, 10000, "", "")
	if err != nil {
		response.InternalError(c, "Failed to export audit logs: "+err.Error())
		return
	}
	response.OK(c, result.Data)
}
