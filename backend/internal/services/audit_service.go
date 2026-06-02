package services

import (
	"encoding/json"
	"fail2ban-dashboard/internal/database"
	"fail2ban-dashboard/internal/models"
	"fmt"
	"math"
	"time"
)

// AuditService handles audit log operations.
type AuditService struct {
	db *database.DB
}

// NewAuditService creates a new AuditService.
func NewAuditService(db *database.DB) *AuditService {
	return &AuditService{db: db}
}

// Log records a new audit entry.
func (s *AuditService) Log(userID int64, username, action, target string, details interface{}, clientIP string) error {
	detailsStr := ""
	if details != nil {
		b, _ := json.Marshal(details)
		detailsStr = string(b)
	}

	_, err := s.db.Exec(
		"INSERT INTO audit_logs (user_id, username, action, target, details, ip_address) VALUES (?, ?, ?, ?, ?, ?)",
		userID, username, action, target, detailsStr, clientIP,
	)
	return err
}

// GetLogs returns paginated audit logs.
func (s *AuditService) GetLogs(page, perPage int, search, action string) (*models.PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	where := "WHERE 1=1"
	args := []interface{}{}
	if search != "" {
		where += " AND (username LIKE ? OR target LIKE ? OR action LIKE ?)"
		pattern := "%" + search + "%"
		args = append(args, pattern, pattern, pattern)
	}
	if action != "" {
		where += " AND action = ?"
		args = append(args, action)
	}

	var total int64
	s.db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM audit_logs %s", where), args...).Scan(&total)

	query := fmt.Sprintf("SELECT id, user_id, username, action, target, details, ip_address, created_at FROM audit_logs %s ORDER BY created_at DESC LIMIT ? OFFSET ?", where)
	args = append(args, perPage, offset)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var l models.AuditLog
		if err := rows.Scan(&l.ID, &l.UserID, &l.Username, &l.Action, &l.Target, &l.Details, &l.IPAddress, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	return &models.PaginatedResponse{
		Data:       logs,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

// GetRecentLogs returns the most recent audit logs.
func (s *AuditService) GetRecentLogs(limit int) ([]models.AuditLog, error) {
	if limit < 1 {
		limit = 10
	}

	rows, err := s.db.Query(
		"SELECT id, user_id, username, action, target, details, ip_address, created_at FROM audit_logs ORDER BY created_at DESC LIMIT ?",
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var l models.AuditLog
		if err := rows.Scan(&l.ID, &l.UserID, &l.Username, &l.Action, &l.Target, &l.Details, &l.IPAddress, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}

// ExportLogs returns all logs in the given time range.
func (s *AuditService) ExportLogs(from, to time.Time) ([]models.AuditLog, error) {
	rows, err := s.db.Query(
		"SELECT id, user_id, username, action, target, details, ip_address, created_at FROM audit_logs WHERE created_at BETWEEN ? AND ? ORDER BY created_at DESC",
		from, to,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var l models.AuditLog
		if err := rows.Scan(&l.ID, &l.UserID, &l.Username, &l.Action, &l.Target, &l.Details, &l.IPAddress, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}
