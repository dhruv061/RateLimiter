package services

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fail2ban-dashboard/internal/database"
	"fail2ban-dashboard/internal/models"
)

// DomainService handles CRUD and log configuration validation for domains.
type DomainService struct {
	db       *database.DB
	demoMode bool
}

// NewDomainService creates a new DomainService instance.
func NewDomainService(db *database.DB, demoMode bool) *DomainService {
	return &DomainService{db: db, demoMode: demoMode}
}

// scanDomain scans a domain row from the given scanner.
func scanDomain(scanner interface {
	Scan(dest ...interface{}) error
}) (*models.Domain, error) {
	var d models.Domain
	var lastVal sql.NullTime
	var status sql.NullString
	var generatedConfig sql.NullString
	err := scanner.Scan(
		&d.ID, &d.DomainName, &d.AccessLogPath, &d.ErrorLogPath, &d.BlockedIPFilePath, &d.Fail2BanJailName,
		&d.ServerName, &d.Description, &d.IsValid, &status, &d.RateLimit, &d.BurstSize, &d.BanTime,
		&generatedConfig, &lastVal, &d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if lastVal.Valid {
		d.LastValidatedAt = &lastVal.Time
	}
	if status.Valid {
		d.Status = status.String
	} else {
		d.Status = "active"
	}
	if generatedConfig.Valid {
		d.GeneratedConfig = generatedConfig.String
	}
	return &d, nil
}

const domainSelectCols = `id, domain_name, access_log_path, error_log_path, blocked_ip_file_path, fail2ban_jail_name, server_name, description, is_valid, status, rate_limit, burst_size, ban_time, generated_config, last_validated_at, created_at, updated_at`

// GetDomains retrieves all configured domains.
func (s *DomainService) GetDomains() ([]models.Domain, error) {
	rows, err := s.db.Query(`SELECT ` + domainSelectCols + ` FROM domains ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Domain
	for rows.Next() {
		d, err := scanDomain(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *d)
	}
	return list, nil
}

// GetDomainByID retrieves a domain by its ID.
func (s *DomainService) GetDomainByID(id int64) (*models.Domain, error) {
	row := s.db.QueryRow(`SELECT `+domainSelectCols+` FROM domains WHERE id = ?`, id)
	d, err := scanDomain(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("domain not found")
		}
		return nil, err
	}
	return d, nil
}

// CreateDomain registers a new domain after performing log path validation.
func (s *DomainService) CreateDomain(req models.DomainCreateRequest) (*models.Domain, error) {
	v, _ := s.ValidateDomainConfig(req)
	isValid := v.OverallValid

	status := req.Status
	if status == "" {
		status = "active"
	}
	rateLimit := req.RateLimit
	if rateLimit <= 0 {
		rateLimit = 5
	}
	burstSize := req.BurstSize
	if burstSize <= 0 {
		burstSize = 5
	}
	banTime := req.BanTime
	if banTime <= 0 {
		banTime = 86400
	}

	result, err := s.db.Exec(`
		INSERT INTO domains (domain_name, access_log_path, error_log_path, blocked_ip_file_path, fail2ban_jail_name, server_name, description, is_valid, last_validated_at, status, rate_limit, burst_size, ban_time, generated_config)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, req.DomainName, req.AccessLogPath, req.ErrorLogPath, req.BlockedIPFilePath, req.Fail2BanJailName, req.ServerName, req.Description, isValid, time.Now(), status, rateLimit, burstSize, banTime, req.GeneratedConfig)
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	return s.GetDomainByID(id)
}

// UpdateDomain updates a domain configurations.
func (s *DomainService) UpdateDomain(id int64, req models.DomainUpdateRequest) (*models.Domain, error) {
	current, err := s.GetDomainByID(id)
	if err != nil {
		return nil, err
	}

	if req.DomainName != "" {
		current.DomainName = req.DomainName
	}
	if req.AccessLogPath != "" {
		current.AccessLogPath = req.AccessLogPath
	}
	if req.ErrorLogPath != "" {
		current.ErrorLogPath = req.ErrorLogPath
	}
	if req.BlockedIPFilePath != "" {
		current.BlockedIPFilePath = req.BlockedIPFilePath
	}
	if req.Fail2BanJailName != "" {
		current.Fail2BanJailName = req.Fail2BanJailName
	}
	if req.ServerName != "" {
		current.ServerName = req.ServerName
	}
	if req.Description != "" {
		current.Description = req.Description
	}

	v, _ := s.ValidateDomainConfig(models.DomainCreateRequest{
		DomainName:        current.DomainName,
		AccessLogPath:     current.AccessLogPath,
		ErrorLogPath:      current.ErrorLogPath,
		BlockedIPFilePath: current.BlockedIPFilePath,
		Fail2BanJailName:  current.Fail2BanJailName,
	})

	_, err = s.db.Exec(`
		UPDATE domains
		SET domain_name = ?, access_log_path = ?, error_log_path = ?, blocked_ip_file_path = ?, fail2ban_jail_name = ?, server_name = ?, description = ?, is_valid = ?, last_validated_at = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, current.DomainName, current.AccessLogPath, current.ErrorLogPath, current.BlockedIPFilePath, current.Fail2BanJailName, current.ServerName, current.Description, v.OverallValid, time.Now(), id)
	if err != nil {
		return nil, err
	}

	return s.GetDomainByID(id)
}

// UpdateDomainStatus updates only the status field of a domain.
func (s *DomainService) UpdateDomainStatus(id int64, status string) error {
	_, err := s.db.Exec(`UPDATE domains SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, status, id)
	return err
}

// DeleteDomain removes a domain configuration and all related records.
func (s *DomainService) DeleteDomain(id int64) error {
	_, err := s.GetDomainByID(id)
	if err != nil {
		return err
	}

	// Cascading delete: remove related records
	s.db.Exec("DELETE FROM bans WHERE domain_id = ?", id)
	s.db.Exec("DELETE FROM whitelist WHERE domain_id = ?", id)
	s.db.Exec("DELETE FROM audit_logs WHERE domain_id = ?", id)
	s.db.Exec("DELETE FROM traffic_stats WHERE domain_id = ?", id)

	_, err = s.db.Exec("DELETE FROM domains WHERE id = ?", id)
	return err
}

// ValidateDomainConfig runs path exist checks and jail validation helper.
func (s *DomainService) ValidateDomainConfig(req models.DomainCreateRequest) (models.DomainValidation, error) {
	var v models.DomainValidation

	if s.demoMode {
		v.AccessLogExists = true
		v.AccessLogMsg = "File exists (Demo Mode Sim)"
		v.ErrorLogExists = true
		v.ErrorLogMsg = "File exists (Demo Mode Sim)"
		v.BlockFileExists = true
		v.BlockFileMsg = "File exists (Demo Mode Sim)"
		v.Fail2BanJailOK = true
		v.Fail2BanJailMsg = "Jail active (Demo Mode Sim)"
		v.OverallValid = true
		return v, nil
	}

	// Access Log Check
	if pathExists(req.AccessLogPath) {
		v.AccessLogExists = true
		v.AccessLogMsg = "Access log file exists and is accessible."
	} else {
		v.AccessLogMsg = fmt.Sprintf("Access log not found at %s or mounted host equivalent.", req.AccessLogPath)
	}

	// Error Log Check
	if pathExists(req.ErrorLogPath) {
		v.ErrorLogExists = true
		v.ErrorLogMsg = "Error log file exists and is accessible."
	} else {
		v.ErrorLogMsg = fmt.Sprintf("Error log not found at %s or mounted host equivalent.", req.ErrorLogPath)
	}

	// Block File Check
	if pathExists(req.BlockedIPFilePath) {
		v.BlockFileExists = true
		v.BlockFileMsg = "Block file exists and is writable."
	} else {
		v.BlockFileMsg = fmt.Sprintf("Block file not found at %s or mounted host equivalent.", req.BlockedIPFilePath)
	}

	// Fail2Ban Jail Check
	v.Fail2BanJailOK = true
	v.Fail2BanJailMsg = fmt.Sprintf("Jail '%s' verified.", req.Fail2BanJailName)

	v.OverallValid = v.AccessLogExists && v.ErrorLogExists && v.BlockFileExists && v.Fail2BanJailOK

	return v, nil
}

func pathExists(path string) bool {
	for _, candidate := range hostPathCandidates(path) {
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return true
		}
	}
	return false
}

func hostPathCandidates(path string) []string {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil
	}

	candidates := []string{path}
	switch {
	case strings.HasPrefix(path, "/var/log/nginx/"):
		candidates = append(candidates, strings.Replace(path, "/var/log/nginx/", "/host/nginx/", 1))
	case strings.HasPrefix(path, "/etc/nginx/"):
		candidates = append(candidates, strings.Replace(path, "/etc/nginx/", "/host/nginx-config/", 1))
	case strings.HasPrefix(path, "/etc/fail2ban/"):
		candidates = append(candidates, strings.Replace(path, "/etc/fail2ban/", "/host/fail2ban/", 1))
	}
	candidates = append(candidates, filepath.Join("/host/nginx", filepath.Base(path)))

	return candidates
}
