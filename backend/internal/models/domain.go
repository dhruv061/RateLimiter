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
	Status            string     `json:"status" db:"status"`                         // pending, active, removing
	RateLimit         int        `json:"rate_limit" db:"rate_limit"`                 // requests per second
	BurstSize         int        `json:"burst_size" db:"burst_size"`                 // burst size
	BanTime           int        `json:"ban_time" db:"ban_time"`                     // ban duration in seconds
	GeneratedConfig   string     `json:"generated_config" db:"generated_config"`     // JSON blob of generated configs
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
	RateLimit         int    `json:"rate_limit"`
	BurstSize         int    `json:"burst_size"`
	BanTime           int    `json:"ban_time"`
	GeneratedConfig   string `json:"generated_config"`
	Status            string `json:"status"`
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
	AccessLogExists bool   `json:"access_log_exists"`
	AccessLogMsg    string `json:"access_log_msg"`
	ErrorLogExists  bool   `json:"error_log_exists"`
	ErrorLogMsg     string `json:"error_log_msg"`
	BlockFileExists bool   `json:"block_file_exists"`
	BlockFileMsg    string `json:"block_file_msg"`
	Fail2BanJailOK  bool   `json:"fail2ban_jail_ok"`
	Fail2BanJailMsg string `json:"fail2ban_jail_msg"`
	OverallValid    bool   `json:"overall_valid"`
}

// GlobalFilter represents the common filter query parameters
// (domain_id, start_time, end_time) parsed from incoming requests.
type GlobalFilter struct {
	DomainID  int64
	StartTime *time.Time
	EndTime   *time.Time
}

// --- Setup Wizard Types ---

// Fail2BanStatusResponse holds the detected Fail2Ban service info.
type Fail2BanStatusResponse struct {
	Installed   bool     `json:"installed"`
	Running     bool     `json:"running"`
	Version     string   `json:"version"`
	ActiveJails []string `json:"active_jails"`
	JailCount   int      `json:"jail_count"`
}

// NginxDiscoveredDomain represents a domain found in Nginx config files.
type NginxDiscoveredDomain struct {
	ServerName string `json:"server_name"`
	ConfigFile string `json:"config_file"`
	HasSSL     bool   `json:"has_ssl"`
}

// NginxDomainDiscoveryResponse returns discovered Nginx domains.
type NginxDomainDiscoveryResponse struct {
	Domains    []NginxDiscoveredDomain `json:"domains"`
	ScannedDirs []string              `json:"scanned_dirs"`
}

// SetupConfigRequest is the payload for generating configuration.
type SetupConfigRequest struct {
	DomainName string `json:"domain_name" binding:"required"`
	RateLimit  int    `json:"rate_limit"`  // default 5
	BurstSize  int    `json:"burst_size"`  // default 5
	BanTime    int    `json:"ban_time"`    // default 86400
}

// GeneratedFile represents a single generated configuration file.
type GeneratedFile struct {
	Filename string `json:"filename"`
	Path     string `json:"path"`
	Content  string `json:"content"`
	Type     string `json:"type"` // filter, action, jail, block, nginx, script
}

// GeneratedConfigResponse contains all generated configuration files.
type GeneratedConfigResponse struct {
	DomainSlug    string          `json:"domain_slug"`
	Files         []GeneratedFile `json:"files"`
	SetupScript   string          `json:"setup_script"`
	NginxSnippet  string          `json:"nginx_snippet"`
	NginxZoneLine string          `json:"nginx_zone_line"`
}

// SetupValidationRequest is the payload for validating a setup.
type SetupValidationRequest struct {
	DomainName string `json:"domain_name" binding:"required"`
}

// SetupValidationCheck represents a single validation check result.
type SetupValidationCheck struct {
	Name    string `json:"name"`
	Passed  bool   `json:"passed"`
	Message string `json:"message"`
}

// SetupValidationResponse contains all validation check results.
type SetupValidationResponse struct {
	Checks       []SetupValidationCheck `json:"checks"`
	OverallValid bool                   `json:"overall_valid"`
}

// CleanupScriptRequest is the payload for generating a cleanup script.
type CleanupScriptRequest struct {
	DomainName string `json:"domain_name" binding:"required"`
}

// CleanupScriptResponse contains the generated cleanup script.
type CleanupScriptResponse struct {
	DomainSlug    string          `json:"domain_slug"`
	CleanupScript string          `json:"cleanup_script"`
	FilesToRemove []GeneratedFile `json:"files_to_remove"`
}

// RemovalValidationRequest is the payload for validating domain removal.
type RemovalValidationRequest struct {
	DomainName string `json:"domain_name" binding:"required"`
}

// RemovalValidationResponse contains removal validation results.
type RemovalValidationResponse struct {
	Checks       []SetupValidationCheck `json:"checks"`
	OverallValid bool                   `json:"overall_valid"`
}
