package handlers

import (
	"fail2ban-dashboard/internal/models"
	"fail2ban-dashboard/internal/services"
	"fail2ban-dashboard/pkg/response"
	"github.com/gin-gonic/gin"
)

// SetupHandler handles REST API calls for the Setup Wizard.
type SetupHandler struct {
	setupSvc *services.SetupService
}

// NewSetupHandler creates a new SetupHandler instance.
func NewSetupHandler(setupSvc *services.SetupService) *SetupHandler {
	return &SetupHandler{setupSvc: setupSvc}
}

// GetFail2BanStatus detects the status of the Fail2Ban service on the host.
func (h *SetupHandler) GetFail2BanStatus(c *gin.Context) {
	status, err := h.setupSvc.GetFail2BanStatus()
	if err != nil {
		response.InternalError(c, "Failed to get Fail2Ban status: "+err.Error())
		return
	}
	response.OK(c, status)
}

// DiscoverNginxDomains scans active Nginx config files for server_name declarations.
func (h *SetupHandler) DiscoverNginxDomains(c *gin.Context) {
	resp, err := h.setupSvc.DiscoverNginxDomains()
	if err != nil {
		response.InternalError(c, "Failed to discover Nginx domains: "+err.Error())
		return
	}
	response.OK(c, resp)
}

// GenerateConfig generates Fail2Ban and Nginx configs for the domain.
func (h *SetupHandler) GenerateConfig(c *gin.Context) {
	var req models.SetupConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid payload: "+err.Error())
		return
	}

	resp, err := h.setupSvc.GenerateConfig(req.DomainName, req.RateLimit, req.BurstSize, req.BanTime)
	if err != nil {
		response.InternalError(c, "Failed to generate config: "+err.Error())
		return
	}
	response.OK(c, resp)
}

// ValidateSetup checks if the generated files and active jails are working on the host.
func (h *SetupHandler) ValidateSetup(c *gin.Context) {
	var req models.SetupValidationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid payload: "+err.Error())
		return
	}

	resp, err := h.setupSvc.ValidateSetup(req.DomainName)
	if err != nil {
		response.InternalError(c, "Failed to validate setup: "+err.Error())
		return
	}
	response.OK(c, resp)
}

// GenerateCleanupScript generates a cleanup shell script for removing domain protection.
func (h *SetupHandler) GenerateCleanupScript(c *gin.Context) {
	var req models.CleanupScriptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid payload: "+err.Error())
		return
	}

	resp, err := h.setupSvc.GenerateCleanupScript(req.DomainName)
	if err != nil {
		response.InternalError(c, "Failed to generate cleanup script: "+err.Error())
		return
	}
	response.OK(c, resp)
}

// ValidateRemoval validates Nginx configs no longer contain domain protection snippets.
func (h *SetupHandler) ValidateRemoval(c *gin.Context) {
	var req models.RemovalValidationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid payload: "+err.Error())
		return
	}

	resp, err := h.setupSvc.ValidateRemoval(req.DomainName)
	if err != nil {
		response.InternalError(c, "Failed to validate removal: "+err.Error())
		return
	}
	response.OK(c, resp)
}
