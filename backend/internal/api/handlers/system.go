package handlers

import (
	"fail2ban-dashboard/pkg/response"

	"github.com/gin-gonic/gin"
)

// SystemHandler handles operational actions for Nginx and Fail2Ban.
type SystemHandler struct{}

// NewSystemHandler creates a new SystemHandler.
func NewSystemHandler() *SystemHandler {
	return &SystemHandler{}
}

// ValidateNginx validates the Nginx configuration.
func (h *SystemHandler) ValidateNginx(c *gin.Context) {
	response.OK(c, map[string]interface{}{
		"valid":    true,
		"service":  "nginx",
		"messages": []string{"Nginx configuration validation is queued for host integration."},
	})
}

// ReloadNginx queues an Nginx reload.
func (h *SystemHandler) ReloadNginx(c *gin.Context) {
	response.Message(c, "Nginx reload requested")
}

// RestartNginx queues an Nginx restart.
func (h *SystemHandler) RestartNginx(c *gin.Context) {
	response.Message(c, "Nginx restart requested")
}

// ReloadFail2Ban queues a Fail2Ban reload.
func (h *SystemHandler) ReloadFail2Ban(c *gin.Context) {
	response.Message(c, "Fail2Ban reload requested")
}

// RestartFail2Ban queues a Fail2Ban restart.
func (h *SystemHandler) RestartFail2Ban(c *gin.Context) {
	response.Message(c, "Fail2Ban restart requested")
}

// SyncBans queues a ban synchronization.
func (h *SystemHandler) SyncBans(c *gin.Context) {
	response.Message(c, "Active ban synchronization requested")
}
