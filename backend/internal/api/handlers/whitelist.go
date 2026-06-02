package handlers

import (
	"strconv"

	"fail2ban-dashboard/internal/models"
	"fail2ban-dashboard/internal/services"
	"fail2ban-dashboard/pkg/response"

	"github.com/gin-gonic/gin"
)

// WhitelistHandler handles whitelist endpoints.
type WhitelistHandler struct {
	wlSvc    *services.WhitelistService
	auditSvc *services.AuditService
}

// NewWhitelistHandler creates a new WhitelistHandler.
func NewWhitelistHandler(wlSvc *services.WhitelistService, auditSvc *services.AuditService) *WhitelistHandler {
	return &WhitelistHandler{wlSvc: wlSvc, auditSvc: auditSvc}
}

// GetAll returns all whitelist entries.
func (h *WhitelistHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "50"))
	search := c.Query("search")

	result, err := h.wlSvc.GetAll(page, perPage, search)
	if err != nil {
		response.InternalError(c, "Failed to fetch whitelist")
		return
	}
	response.OK(c, result)
}

// Add creates a new whitelist entry.
func (h *WhitelistHandler) Add(c *gin.Context) {
	var req models.WhitelistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: ip_address required")
		return
	}

	username, _ := c.Get("username")
	entry, err := h.wlSvc.Add(req.IPAddress, req.Description, username.(string))
	if err != nil {
		response.BadRequest(c, "Failed to add whitelist entry: "+err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	h.auditSvc.Log(userID.(int64), username.(string), "add_whitelist", req.IPAddress, nil, c.ClientIP())

	response.Created(c, entry)
}

// Remove deletes a whitelist entry.
func (h *WhitelistHandler) Remove(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid whitelist ID")
		return
	}

	if err := h.wlSvc.Remove(id); err != nil {
		response.NotFound(c, "Whitelist entry not found")
		return
	}

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	h.auditSvc.Log(userID.(int64), username.(string), "remove_whitelist", "", map[string]int64{"id": id}, c.ClientIP())

	response.Message(c, "Whitelist entry removed")
}

// Export returns all whitelist entries for export.
func (h *WhitelistHandler) Export(c *gin.Context) {
	result, err := h.wlSvc.GetAll(1, 10000, "")
	if err != nil {
		response.InternalError(c, "Failed to export whitelist")
		return
	}
	response.OK(c, result.Data)
}
