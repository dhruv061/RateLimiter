package handlers

import (
	"strconv"

	"fail2ban-dashboard/internal/models"
	"fail2ban-dashboard/internal/services"
	"fail2ban-dashboard/pkg/response"

	"github.com/gin-gonic/gin"
)

// BansHandler handles ban-related endpoints.
type BansHandler struct {
	banSvc   *services.BanService
	auditSvc *services.AuditService
}

// NewBansHandler creates a new BansHandler.
func NewBansHandler(banSvc *services.BanService, auditSvc *services.AuditService) *BansHandler {
	return &BansHandler{banSvc: banSvc, auditSvc: auditSvc}
}

// GetActiveBans returns paginated active bans scoped by global filters.
func (h *BansHandler) GetActiveBans(c *gin.Context) {
	filter := parseGlobalFilter(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	search := c.Query("search")
	sortBy := c.DefaultQuery("sort_by", "ban_time")
	sortDir := c.DefaultQuery("sort_dir", "desc")
	country := c.Query("country")

	result, err := h.banSvc.GetActiveBans(filter, page, perPage, search, sortBy, sortDir, country)
	if err != nil {
		response.InternalError(c, "Failed to fetch active bans: "+err.Error())
		return
	}
	response.OK(c, result)
}

// GetBanHistory returns paginated ban history scoped by global filters.
func (h *BansHandler) GetBanHistory(c *gin.Context) {
	filter := parseGlobalFilter(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	search := c.Query("search")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	result, err := h.banSvc.GetBanHistory(filter, page, perPage, search, dateFrom, dateTo)
	if err != nil {
		response.InternalError(c, "Failed to fetch ban history: "+err.Error())
		return
	}
	response.OK(c, result)
}

// GetBanDetail returns a single ban's full details.
func (h *BansHandler) GetBanDetail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid ban ID")
		return
	}

	ban, err := h.banSvc.GetBanByID(id)
	if err != nil {
		response.NotFound(c, "Ban not found: "+err.Error())
		return
	}
	response.OK(c, ban)
}

// UnbanIP marks a ban as inactive.
func (h *BansHandler) UnbanIP(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid ban ID")
		return
	}

	// Get ban details for audit
	ban, err := h.banSvc.GetBanByID(id)
	if err != nil {
		response.NotFound(c, "Ban not found")
		return
	}

	if err := h.banSvc.UnbanIP(id); err != nil {
		response.InternalError(c, "Failed to unban IP")
		return
	}

	// Audit log
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	h.auditSvc.Log(ban.DomainID, userID.(int64), username.(string), "unban_ip", ban.IPAddress, map[string]interface{}{
		"ban_id": id,
		"jail":   ban.Jail,
	}, c.ClientIP())

	response.Message(c, "IP "+ban.IPAddress+" successfully unbanned")
}

// BulkUnban marks multiple bans as inactive.
func (h *BansHandler) BulkUnban(c *gin.Context) {
	var req models.BulkUnbanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: ids array required")
		return
	}

	filter := parseGlobalFilter(c)
	affected, err := h.banSvc.BulkUnban(req.IDs)
	if err != nil {
		response.InternalError(c, "Failed to perform bulk unban")
		return
	}

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	h.auditSvc.Log(filter.DomainID, userID.(int64), username.(string), "bulk_unban", "", map[string]interface{}{
		"count": affected,
		"ids":   req.IDs,
	}, c.ClientIP())

	response.OK(c, map[string]int64{"unbanned": affected})
}

// GetTopOffenders returns the most frequently banned IPs.
func (h *BansHandler) GetTopOffenders(c *gin.Context) {
	filter := parseGlobalFilter(c)
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	offenders, err := h.banSvc.GetTopOffenders(filter, limit)
	if err != nil {
		response.InternalError(c, "Failed to fetch top offenders: "+err.Error())
		return
	}
	response.OK(c, offenders)
}
