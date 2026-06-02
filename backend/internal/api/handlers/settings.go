package handlers

import (
	"fail2ban-dashboard/internal/models"
	"fail2ban-dashboard/internal/services"
	"fail2ban-dashboard/pkg/response"

	"github.com/gin-gonic/gin"
)

// SettingsHandler handles settings endpoints.
type SettingsHandler struct {
	settingsSvc *services.SettingsService
	auditSvc    *services.AuditService
}

// NewSettingsHandler creates a new SettingsHandler.
func NewSettingsHandler(settingsSvc *services.SettingsService, auditSvc *services.AuditService) *SettingsHandler {
	return &SettingsHandler{settingsSvc: settingsSvc, auditSvc: auditSvc}
}

// GetSettings returns all settings.
func (h *SettingsHandler) GetSettings(c *gin.Context) {
	settings, err := h.settingsSvc.GetAll()
	if err != nil {
		response.InternalError(c, "Failed to fetch settings")
		return
	}
	response.OK(c, settings)
}

// UpdateSettings updates multiple settings.
func (h *SettingsHandler) UpdateSettings(c *gin.Context) {
	var req models.SettingsUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: settings map required")
		return
	}

	username, _ := c.Get("username")
	if err := h.settingsSvc.Update(req.Settings, username.(string)); err != nil {
		response.InternalError(c, "Failed to update settings")
		return
	}

	userID, _ := c.Get("user_id")
	h.auditSvc.Log(userID.(int64), username.(string), "update_settings", "", req.Settings, c.ClientIP())

	response.Message(c, "Settings updated successfully")
}

// ValidateConfig validates the current configuration (placeholder).
func (h *SettingsHandler) ValidateConfig(c *gin.Context) {
	response.OK(c, map[string]interface{}{
		"valid":    true,
		"messages": []string{},
	})
}
