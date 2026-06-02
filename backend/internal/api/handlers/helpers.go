package handlers

import (
	"strconv"
	"time"

	"fail2ban-dashboard/internal/models"
	"github.com/gin-gonic/gin"
)

// parseGlobalFilter extracts domain_id, start_time, and end_time query parameters.
func parseGlobalFilter(c *gin.Context) models.GlobalFilter {
	var filter models.GlobalFilter

	// Parse Domain ID
	domainIDStr := c.Query("domain_id")
	if domainIDStr != "" {
		if id, err := strconv.ParseInt(domainIDStr, 10, 64); err == nil {
			filter.DomainID = id
		}
	}

	// Parse Start Time (accept RFC3339, ISO8601 with milliseconds, or simple Date)
	startTimeStr := c.Query("start_time")
	if startTimeStr != "" {
		if t, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
			filter.StartTime = &t
		} else if t, err := time.Parse("2006-01-02T15:04:05.000Z", startTimeStr); err == nil {
			filter.StartTime = &t
		} else if t, err := time.Parse("2006-01-02", startTimeStr); err == nil {
			filter.StartTime = &t
		}
	}

	// Parse End Time
	endTimeStr := c.Query("end_time")
	if endTimeStr != "" {
		if t, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
			filter.EndTime = &t
		} else if t, err := time.Parse("2006-01-02T15:04:05.000Z", endTimeStr); err == nil {
			filter.EndTime = &t
		} else if t, err := time.Parse("2006-01-02", endTimeStr); err == nil {
			filter.EndTime = &t
		}
	}

	return filter
}
