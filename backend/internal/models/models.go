package models

import "time"

// User represents a dashboard user.
type User struct {
	ID              int64      `json:"id" db:"id"`
	Username        string     `json:"username" db:"username"`
	PasswordHash    string     `json:"-" db:"password_hash"`
	Email           string     `json:"email" db:"email"`
	Role            string     `json:"role" db:"role"`
	MustChangePass  bool       `json:"must_change_pass" db:"must_change_pass"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
	LastLogin       *time.Time `json:"last_login,omitempty" db:"last_login"`
}

// Ban represents an IP ban record.
type Ban struct {
	ID             int64      `json:"id" db:"id"`
	IPAddress      string     `json:"ip_address" db:"ip_address"`
	Country        string     `json:"country" db:"country"`
	CountryCode    string     `json:"country_code" db:"country_code"`
	Region         string     `json:"region" db:"region"`
	City           string     `json:"city" db:"city"`
	ASN            string     `json:"asn" db:"asn"`
	ISP            string     `json:"isp" db:"isp"`
	Jail           string     `json:"jail" db:"jail"`
	Reason         string     `json:"reason" db:"reason"`
	BanTime        time.Time  `json:"ban_time" db:"ban_time"`
	UnbanTime      *time.Time `json:"unban_time,omitempty" db:"unban_time"`
	BanDuration    int        `json:"ban_duration" db:"ban_duration"`
	RequestCount   int        `json:"request_count" db:"request_count"`
	ViolationCount int        `json:"violation_count" db:"violation_count"`
	IsActive       bool       `json:"is_active" db:"is_active"`
	DomainID       int64      `json:"domain_id" db:"domain_id"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
}

// RemainingTime calculates seconds remaining on an active ban.
func (b *Ban) RemainingTime() int {
	if !b.IsActive {
		return 0
	}
	end := b.BanTime.Add(time.Duration(b.BanDuration) * time.Second)
	remaining := time.Until(end)
	if remaining < 0 {
		return 0
	}
	return int(remaining.Seconds())
}

// WhitelistEntry represents a whitelisted IP.
type WhitelistEntry struct {
	ID          int64     `json:"id" db:"id"`
	IPAddress   string    `json:"ip_address" db:"ip_address"`
	Description string    `json:"description" db:"description"`
	AddedBy     string    `json:"added_by" db:"added_by"`
	DomainID    int64     `json:"domain_id" db:"domain_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// AuditLog represents a dashboard action log entry.
type AuditLog struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Username  string    `json:"username" db:"username"`
	Action    string    `json:"action" db:"action"`
	Target    string    `json:"target" db:"target"`
	Details   string    `json:"details" db:"details"`
	IPAddress string    `json:"ip_address" db:"ip_address"`
	DomainID  int64     `json:"domain_id" db:"domain_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Setting represents a key-value configuration entry.
type Setting struct {
	ID        int64     `json:"id" db:"id"`
	Key       string    `json:"key" db:"key"`
	Value     string    `json:"value" db:"value"`
	UpdatedBy string    `json:"updated_by" db:"updated_by"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// TrafficStat represents aggregated traffic metrics for a time bucket.
type TrafficStat struct {
	ID              int64     `json:"id" db:"id"`
	Timestamp       time.Time `json:"timestamp" db:"timestamp"`
	TotalRequests   int       `json:"total_requests" db:"total_requests"`
	UniqueIPs       int       `json:"unique_ips" db:"unique_ips"`
	Status429       int       `json:"status_429" db:"status_429"`
	Status403       int       `json:"status_403" db:"status_403"`
	AvgResponseTime float64   `json:"avg_response_time" db:"avg_response_time"`
	Period          string    `json:"period" db:"period"`
	DomainID        int64     `json:"domain_id" db:"domain_id"`
}

// LiveRequest represents a single parsed Nginx access log entry.
type LiveRequest struct {
	Timestamp    time.Time `json:"timestamp"`
	IPAddress    string    `json:"ip_address"`
	Method       string    `json:"method"`
	URL          string    `json:"url"`
	StatusCode   int       `json:"status_code"`
	ResponseTime float64   `json:"response_time"`
	UserAgent    string    `json:"user_agent"`
	BytesSent    int64     `json:"bytes_sent"`
}

// DashboardStats holds the overview metrics.
type DashboardStats struct {
	// Traffic
	TotalRequestsToday  int     `json:"total_requests_today"`
	RequestsLastHour    int     `json:"requests_last_hour"`
	CurrentRPS          float64 `json:"current_rps"`
	PeakRPS             float64 `json:"peak_rps"`

	// Rate Limiting
	Total429Today    int     `json:"total_429_today"`
	Count429LastHour int     `json:"count_429_last_hour"`
	Current429Rate   float64 `json:"current_429_rate"`

	// Bans
	TotalBansToday int `json:"total_bans_today"`
	ActiveBans     int `json:"active_bans"`
	Bans24h        int `json:"bans_24h"`
	UnbansToday    int `json:"unbans_today"`

	// System
	NginxStatus    string `json:"nginx_status"`
	Fail2BanStatus string `json:"fail2ban_status"`
	DatabaseStatus string `json:"database_status"`
	ServiceStatus  string `json:"service_status"`
}

// SystemStatus holds service health info.
type SystemStatus struct {
	Nginx    ServiceInfo `json:"nginx"`
	Fail2Ban ServiceInfo `json:"fail2ban"`
	Database ServiceInfo `json:"database"`
}

// ServiceInfo describes a service's health.
type ServiceInfo struct {
	Status  string `json:"status"` // "running", "stopped", "unknown"
	Version string `json:"version,omitempty"`
	Uptime  string `json:"uptime,omitempty"`
	Details string `json:"details,omitempty"`
}

// AttackStatus describes the current threat level.
type AttackStatus struct {
	Level           string `json:"level"` // "normal", "elevated", "attack"
	RequestsLastMin int    `json:"requests_last_min"`
	UniqueIPsLastMin int   `json:"unique_ips_last_min"`
	Count429LastMin int    `json:"count_429_last_min"`
	BansLastMin     int    `json:"bans_last_min"`
	RequestsLast5Min int   `json:"requests_last_5min"`
	UniqueIPsLast5Min int  `json:"unique_ips_last_5min"`
	Count429Last5Min int   `json:"count_429_last_5min"`
	BansLast5Min     int   `json:"bans_last_5min"`
}

// TopOffender represents an IP with high violation counts.
type TopOffender struct {
	IPAddress      string `json:"ip_address"`
	Country        string `json:"country"`
	CountryCode    string `json:"country_code"`
	TotalRequests  int    `json:"total_requests"`
	ViolationCount int    `json:"violation_count"`
	BanCount       int    `json:"ban_count"`
}

// CountryStats aggregates analytics per country.
type CountryStats struct {
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	Requests    int    `json:"requests"`
	Violations  int    `json:"violations"`
	Bans        int    `json:"bans"`
}

// PaginatedResponse wraps paginated list results.
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	TotalPages int         `json:"total_pages"`
}

// LoginRequest is the login payload.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse is the login response.
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// WhitelistRequest is the payload for adding a whitelist entry.
type WhitelistRequest struct {
	IPAddress   string `json:"ip_address" binding:"required"`
	Description string `json:"description"`
	DomainID    int64  `json:"domain_id"`
}

// BulkUnbanRequest is the payload for bulk unban.
type BulkUnbanRequest struct {
	IDs []int64 `json:"ids" binding:"required"`
}

// SettingsUpdateRequest is the payload for updating settings.
type SettingsUpdateRequest struct {
	Settings map[string]string `json:"settings" binding:"required"`
}
