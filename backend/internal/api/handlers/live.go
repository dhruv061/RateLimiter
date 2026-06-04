package handlers

import (
	"bytes"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"fail2ban-dashboard/internal/models"
	"fail2ban-dashboard/pkg/response"

	"github.com/gin-gonic/gin"
)

// LiveHandler handles live request endpoints.
type LiveHandler struct {
	accessLogPath string
}

// NewLiveHandler creates a new LiveHandler.
func NewLiveHandler(accessLogPath string) *LiveHandler {
	return &LiveHandler{accessLogPath: accessLogPath}
}

var accessLogPattern = regexp.MustCompile(`^(\S+) \S+ \S+ \[([^\]]+)\] "([^"]*)" (\d{3}) (\S+) "[^"]*" "([^"]*)"`)

// GetRequests returns recent live requests parsed from the configured Nginx access log.
func (h *LiveHandler) GetRequests(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if limit < 1 || limit > 200 {
		limit = 50
	}

	lines, err := tailLines(h.accessLogPath, limit*3)
	if err != nil {
		response.OK(c, []models.LiveRequest{})
		return
	}

	requests := make([]models.LiveRequest, 0, limit)
	for i := len(lines) - 1; i >= 0 && len(requests) < limit; i-- {
		if req, ok := parseAccessLogLine(lines[i]); ok {
			requests = append(requests, req)
		}
	}

	response.OK(c, requests)
}

func tailLines(path string, limit int) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	const maxTailBytes int64 = 1024 * 1024
	offset := info.Size() - maxTailBytes
	if offset < 0 {
		offset = 0
	}
	if _, err := file.Seek(offset, io.SeekStart); err != nil {
		return nil, err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if offset > 0 {
		if newline := bytes.IndexByte(data, '\n'); newline >= 0 {
			data = data[newline+1:]
		}
	}

	all := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	if len(all) > limit {
		all = all[len(all)-limit:]
	}
	return all, nil
}

func parseAccessLogLine(line string) (models.LiveRequest, bool) {
	matches := accessLogPattern.FindStringSubmatch(line)
	if matches == nil {
		return models.LiveRequest{}, false
	}

	timestamp, err := time.Parse("02/Jan/2006:15:04:05 -0700", matches[2])
	if err != nil {
		return models.LiveRequest{}, false
	}

	requestParts := strings.Fields(matches[3])
	method := ""
	requestURL := matches[3]
	if len(requestParts) >= 2 {
		method = requestParts[0]
		requestURL = requestParts[1]
	}

	status, _ := strconv.Atoi(matches[4])
	bytesSent, _ := strconv.ParseInt(matches[5], 10, 64)

	return models.LiveRequest{
		Timestamp:    timestamp,
		IPAddress:    matches[1],
		Method:       method,
		URL:          requestURL,
		StatusCode:   status,
		ResponseTime: 0,
		UserAgent:    matches[6],
		BytesSent:    bytesSent,
	}, true
}
