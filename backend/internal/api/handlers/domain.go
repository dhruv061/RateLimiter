package handlers

import (
	"strconv"

	"fail2ban-dashboard/internal/models"
	"fail2ban-dashboard/internal/services"
	"fail2ban-dashboard/pkg/response"
	"github.com/gin-gonic/gin"
)

// DomainHandler exposes REST handlers for domains.
type DomainHandler struct {
	domainSvc *services.DomainService
}

// NewDomainHandler creates a new DomainHandler instance.
func NewDomainHandler(domainSvc *services.DomainService) *DomainHandler {
	return &DomainHandler{domainSvc: domainSvc}
}

// GetDomains returns a list of all configured domains.
func (h *DomainHandler) GetDomains(c *gin.Context) {
	list, err := h.domainSvc.GetDomains()
	if err != nil {
		response.InternalError(c, "Failed to fetch domains: "+err.Error())
		return
	}
	response.OK(c, list)
}

// CreateDomain parses the payload and adds a new domain configuration.
func (h *DomainHandler) CreateDomain(c *gin.Context) {
	var req models.DomainCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request payload: "+err.Error())
		return
	}

	d, err := h.domainSvc.CreateDomain(req)
	if err != nil {
		response.InternalError(c, "Failed to create domain: "+err.Error())
		return
	}
	response.Created(c, d)
}

// DeleteDomain deletes a domain configuration.
func (h *DomainHandler) DeleteDomain(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid domain ID")
		return
	}

	err = h.domainSvc.DeleteDomain(id)
	if err != nil {
		response.InternalError(c, "Failed to delete domain: "+err.Error())
		return
	}
	response.OK(c, gin.H{"message": "Domain deleted successfully"})
}

// ValidateDomain runs log check validations on a payload before saving.
func (h *DomainHandler) ValidateDomain(c *gin.Context) {
	var req models.DomainCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request payload: "+err.Error())
		return
	}

	v, err := h.domainSvc.ValidateDomainConfig(req)
	if err != nil {
		response.InternalError(c, "Failed to validate configuration: "+err.Error())
		return
	}
	response.OK(c, v)
}
