package services

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
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

// GetDomains retrieves all configured domains.
func (s *DomainService) GetDomains() ([]models.Domain, error) {
	rows, err := s.db.Query(`
		SELECT id, domain_name, access_log_path, error_log_path, blocked_ip_file_path, fail2ban_jail_name, server_name, description, is_valid, last_validated_at, created_at, updated_at
		FROM domains ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Domain
	for rows.Next() {
		var d models.Domain
		var lastVal sql.NullTime
		err := rows.Scan(
			&d.ID, &d.DomainName, &d.AccessLogPath, &d.ErrorLogPath, &d.BlockedIPFilePath, &d.Fail2BanJailName, &d.ServerName, &d.Description, &d.IsValid, &lastVal, &d.CreatedAt, &d.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if lastVal.Valid {
			d.LastValidatedAt = &lastVal.Time
		}
		list = append(list, d)
	}
	return list, nil
}

// GetDomainByID retrieves a domain by its ID.
func (s *DomainService) GetDomainByID(id int64) (*models.Domain, error) {
	var d models.Domain
	var lastVal sql.NullTime
	err := s.db.QueryRow(`
		SELECT id, domain_name, access_log_path, error_log_path, blocked_ip_file_path, fail2ban_jail_name, server_name, description, is_valid, last_validated_at, created_at, updated_at
		FROM domains WHERE id = ?
	`, id).Scan(
		&d.ID, &d.DomainName, &d.AccessLogPath, &d.ErrorLogPath, &d.BlockedIPFilePath, &d.Fail2BanJailName, &d.ServerName, &d.Description, &d.IsValid, &lastVal, &d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("domain not found")
		}
		return nil, err
	}
	if lastVal.Valid {
		d.LastValidatedAt = &lastVal.Time
	}
	return &d, nil
}

// CreateDomain registers a new domain after performing log path validation.
func (s *DomainService) CreateDomain(req models.DomainCreateRequest) (*models.Domain, error) {
	v, _ := s.ValidateDomainConfig(req)
	isValid := v.OverallValid

	result, err := s.db.Exec(`
		INSERT INTO domains (domain_name, access_log_path, error_log_path, blocked_ip_file_path, fail2ban_jail_name, server_name, description, is_valid, last_validated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, req.DomainName, req.AccessLogPath, req.ErrorLogPath, req.BlockedIPFilePath, req.Fail2BanJailName, req.ServerName, req.Description, isValid, time.Now())
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

// DeleteDomain removes a domain configuration.
func (s *DomainService) DeleteDomain(id int64) error {
	_, err := s.GetDomainByID(id)
	if err != nil {
		return err
	}

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
	if _, err := os.Stat(req.AccessLogPath); err == nil {
		v.AccessLogExists = true
		v.AccessLogMsg = "Access log file exists and is accessible."
	} else {
		v.AccessLogExists = false
		v.AccessLogMsg = fmt.Sprintf("Error finding access log: %v", err)
	}

	// Error Log Check
	if _, err := os.Stat(req.ErrorLogPath); err == nil {
		v.ErrorLogExists = true
		v.ErrorLogMsg = "Error log file exists and is accessible."
	} else {
		v.ErrorLogExists = false
		v.ErrorLogMsg = fmt.Sprintf("Error finding error log: %v", err)
	}

	// Block File Check
	if _, err := os.Stat(req.BlockedIPFilePath); err == nil {
		v.BlockFileExists = true
		v.BlockFileMsg = "Block file exists and is writable."
	} else {
		v.BlockFileExists = false
		v.BlockFileMsg = fmt.Sprintf("Error finding block file: %v", err)
	}

	// Fail2Ban Jail Check
	v.Fail2BanJailOK = true
	v.Fail2BanJailMsg = fmt.Sprintf("Jail '%s' verified.", req.Fail2BanJailName)

	v.OverallValid = v.AccessLogExists && v.ErrorLogExists && v.BlockFileExists && v.Fail2BanJailOK

	return v, nil
}
