package handlers

import (
	"strconv"
	"time"

	"fail2ban-dashboard/internal/models"
	"fail2ban-dashboard/pkg/response"

	"github.com/gin-gonic/gin"
)

// LiveHandler handles live request endpoints.
type LiveHandler struct{}

// NewLiveHandler creates a new LiveHandler.
func NewLiveHandler() *LiveHandler {
	return &LiveHandler{}
}

// GetRequests returns recent live requests. Real log tailing can feed this endpoint later.
func (h *LiveHandler) GetRequests(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if limit < 1 || limit > 200 {
		limit = 50
	}

	requests := make([]models.LiveRequest, 0, limit)
	statuses := []int{200, 200, 200, 301, 403, 429}
	methods := []string{"GET", "POST", "GET", "GET", "POST"}
	urls := []string{"/", "/api/auth/login", "/wp-login.php", "/api/dashboard/stats", "/.env", "/admin"}
	agents := []string{"Mozilla/5.0", "curl/8.4.0", "python-requests/2.31.0", "Go-http-client/2.0"}

	now := time.Now()
	for i := 0; i < limit; i++ {
		requests = append(requests, models.LiveRequest{
			Timestamp:    now.Add(-time.Duration(i*7) * time.Second),
			IPAddress:    "203.0.113." + strconv.Itoa((i%200)+1),
			Method:       methods[i%len(methods)],
			URL:          urls[i%len(urls)],
			StatusCode:   statuses[i%len(statuses)],
			ResponseTime: float64(18 + (i % 120)),
			UserAgent:    agents[i%len(agents)],
			BytesSent:    int64(512 + i*37),
		})
	}

	response.OK(c, requests)
}
