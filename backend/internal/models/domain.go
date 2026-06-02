package models

import "time"

// Domain represents a managed domain configuration.
type Domain struct {
	ID                int64      `json:"id" db:"id"`
	DomainName        string     `json:"domain_name" db:"domain_name"`
	AccessLogPath     string     `json:"access_log_path" db:"access_log_path"`
	ErrorLogPath      string     `json:"error_log_path" db:"error_log_path"`
	BlockedIPFilePath string     `json:"blocked_ip_file_path" db:"blocked_ip_file_path"`
	Fail2BanJailName  string     `json:"fail2ban_jail_name" db:"fail2ban_jail_name"`
	ServerName        string     `json:"server_name" db:"server_name"`
	Description       string     `json:"description" db:"description"`
	IsValid           bool       `json:"is_valid" db:"is_valid"`
	LastValidatedAt   *time.Time `json:"last_validated_at,omitempty" db:"last_validated_at"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

// DomainCreateRequest is the payload for adding a new domain.
type DomainCreateRequest struct {
	DomainName        string `json:"domain_name" binding:"required"`
	AccessLogPath     string `json:"access_log_path" binding:"required"`
	ErrorLogPath      string `json:"error_log_path" binding:"required"`
	BlockedIPFilePath string `json:"blocked_ip_file_path" binding:"required"`
	Fail2BanJailName  string `json:"fail2ban_jail_name" binding:"required"`
	ServerName        string `json:"server_name"`
	Description       string `json:"description"`
}

// DomainUpdateRequest is the payload for updating a domain.
type DomainUpdateRequest struct {
	DomainName        string `json:"domain_name"`
	AccessLogPath     string `json:"access_log_path"`
	ErrorLogPath      string `json:"error_log_path"`
	BlockedIPFilePath string `json:"blocked_ip_file_path"`
	Fail2BanJailName  string `json:"fail2ban_jail_name"`
	ServerName        string `json:"server_name"`
	Description       string `json:"description"`
}

// DomainValidation holds the results of a domain config validation.
type DomainValidation struct {
	AccessLogExists  bool   `json:"access_log_exists"`
	AccessLogMsg     string `json:"access_log_msg"`
	ErrorLogExists   bool   `json:"error_log_exists"`
	ErrorLogMsg      string `json:"error_log_msg"`
	BlockFileExists  bool   `json:"block_file_exists"`
	BlockFileMsg     string `json:"block_file_msg"`
	Fail2BanJailOK   bool   `json:"fail2ban_jail_ok"`
	Fail2BanJailMsg  string `json:"fail2ban_jail_msg"`
	OverallValid     bool   `json:"overall_valid"`
}

// GlobalFilter represents the common filter query parameters
// (domain_id, start_time, end_time) parsed from incoming requests.
type GlobalFilter struct {
	DomainID  int64
	StartTime *time.Time
	EndTime   *time.Time
}
