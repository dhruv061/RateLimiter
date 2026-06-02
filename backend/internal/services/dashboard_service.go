package services

import (
	"fmt"
	"time"

	"fail2ban-dashboard/internal/database"
	"fail2ban-dashboard/internal/models"
)

// DashboardService aggregates dashboard metrics.
type DashboardService struct {
	db *database.DB
}

// NewDashboardService creates a new DashboardService.
func NewDashboardService(db *database.DB) *DashboardService {
	return &DashboardService{db: db}
}

// GetStats returns aggregated dashboard statistics.
func (s *DashboardService) GetStats() (*models.DashboardStats, error) {
	stats := &models.DashboardStats{}
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	hourAgo := now.Add(-1 * time.Hour)
	dayAgo := now.Add(-24 * time.Hour)

	// Traffic stats from traffic_stats table
	s.db.QueryRow(
		"SELECT COALESCE(SUM(total_requests), 0) FROM traffic_stats WHERE timestamp >= ? AND period = 'hour'",
		todayStart,
	).Scan(&stats.TotalRequestsToday)

	s.db.QueryRow(
		"SELECT COALESCE(SUM(total_requests), 0) FROM traffic_stats WHERE timestamp >= ? AND period = 'minute'",
		hourAgo,
	).Scan(&stats.RequestsLastHour)

	// Rate limiting
	s.db.QueryRow(
		"SELECT COALESCE(SUM(status_429), 0) FROM traffic_stats WHERE timestamp >= ? AND period = 'hour'",
		todayStart,
	).Scan(&stats.Total429Today)

	s.db.QueryRow(
		"SELECT COALESCE(SUM(status_429), 0) FROM traffic_stats WHERE timestamp >= ? AND period = 'minute'",
		hourAgo,
	).Scan(&stats.Count429LastHour)

	// Ban stats
	s.db.QueryRow(
		"SELECT COUNT(*) FROM bans WHERE ban_time >= ?",
		todayStart,
	).Scan(&stats.TotalBansToday)

	s.db.QueryRow(
		"SELECT COUNT(*) FROM bans WHERE is_active = 1",
	).Scan(&stats.ActiveBans)

	s.db.QueryRow(
		"SELECT COUNT(*) FROM bans WHERE ban_time >= ?",
		dayAgo,
	).Scan(&stats.Bans24h)

	s.db.QueryRow(
		"SELECT COUNT(*) FROM bans WHERE unban_time IS NOT NULL AND unban_time >= ?",
		todayStart,
	).Scan(&stats.UnbansToday)

	// System status (will be updated by system service)
	stats.NginxStatus = "running"
	stats.Fail2BanStatus = "running"
	stats.DatabaseStatus = "running"
	stats.ServiceStatus = "running"

	return stats, nil
}

// GetSystemStatus returns service health information.
func (s *DashboardService) GetSystemStatus() *models.SystemStatus {
	return &models.SystemStatus{
		Nginx: models.ServiceInfo{
			Status:  "running",
			Version: "1.24.0",
			Details: "Active",
		},
		Fail2Ban: models.ServiceInfo{
			Status:  "running",
			Version: "0.11.2",
			Details: "Active",
		},
		Database: models.ServiceInfo{
			Status:  "running",
			Version: "SQLite 3.45",
			Details: "WAL mode",
		},
	}
}

// GetAttackStatus evaluates the current threat level.
func (s *DashboardService) GetAttackStatus() (*models.AttackStatus, error) {
	status := &models.AttackStatus{}
	now := time.Now()
	oneMinAgo := now.Add(-1 * time.Minute)
	fiveMinAgo := now.Add(-5 * time.Minute)

	// Last minute stats
	s.db.QueryRow(
		"SELECT COALESCE(SUM(total_requests), 0) FROM traffic_stats WHERE timestamp >= ? AND period = 'minute'",
		oneMinAgo,
	).Scan(&status.RequestsLastMin)

	s.db.QueryRow(
		"SELECT COALESCE(SUM(unique_ips), 0) FROM traffic_stats WHERE timestamp >= ? AND period = 'minute'",
		oneMinAgo,
	).Scan(&status.UniqueIPsLastMin)

	s.db.QueryRow(
		"SELECT COALESCE(SUM(status_429), 0) FROM traffic_stats WHERE timestamp >= ? AND period = 'minute'",
		oneMinAgo,
	).Scan(&status.Count429LastMin)

	s.db.QueryRow(
		"SELECT COUNT(*) FROM bans WHERE ban_time >= ?",
		oneMinAgo,
	).Scan(&status.BansLastMin)

	// Last 5 minutes stats
	s.db.QueryRow(
		"SELECT COALESCE(SUM(total_requests), 0) FROM traffic_stats WHERE timestamp >= ? AND period = 'minute'",
		fiveMinAgo,
	).Scan(&status.RequestsLast5Min)

	s.db.QueryRow(
		"SELECT COALESCE(SUM(unique_ips), 0) FROM traffic_stats WHERE timestamp >= ? AND period = 'minute'",
		fiveMinAgo,
	).Scan(&status.UniqueIPsLast5Min)

	s.db.QueryRow(
		"SELECT COALESCE(SUM(status_429), 0) FROM traffic_stats WHERE timestamp >= ? AND period = 'minute'",
		fiveMinAgo,
	).Scan(&status.Count429Last5Min)

	s.db.QueryRow(
		"SELECT COUNT(*) FROM bans WHERE ban_time >= ?",
		fiveMinAgo,
	).Scan(&status.BansLast5Min)

	// Determine threat level
	status.Level = "normal"
	if status.Count429LastMin > 10 || status.BansLastMin > 3 {
		status.Level = "elevated"
	}
	if status.Count429LastMin > 50 || status.BansLastMin > 10 {
		status.Level = "attack"
	}

	return status, nil
}

// GetTrafficTrends returns traffic data for charts.
func (s *DashboardService) GetTrafficTrends(period string, hours int) ([]models.TrafficStat, error) {
	if hours < 1 {
		hours = 24
	}
	since := time.Now().Add(-time.Duration(hours) * time.Hour)

	rows, err := s.db.Query(
		"SELECT id, timestamp, total_requests, unique_ips, status_429, status_403, avg_response_time, period FROM traffic_stats WHERE period = ? AND timestamp >= ? ORDER BY timestamp ASC",
		period, since,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query traffic trends: %w", err)
	}
	defer rows.Close()

	var stats []models.TrafficStat
	for rows.Next() {
		var ts models.TrafficStat
		if err := rows.Scan(&ts.ID, &ts.Timestamp, &ts.TotalRequests, &ts.UniqueIPs, &ts.Status429, &ts.Status403, &ts.AvgResponseTime, &ts.Period); err != nil {
			return nil, err
		}
		stats = append(stats, ts)
	}
	return stats, nil
}

// GetCountryStats returns analytics grouped by country.
func (s *DashboardService) GetCountryStats() ([]models.CountryStats, error) {
	rows, err := s.db.Query(`
		SELECT country, country_code,
			COUNT(*) as ban_count,
			COALESCE(SUM(request_count), 0) as total_requests,
			COALESCE(SUM(violation_count), 0) as total_violations
		FROM bans
		WHERE country != ''
		GROUP BY country_code
		ORDER BY ban_count DESC
		LIMIT 50
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []models.CountryStats
	for rows.Next() {
		var cs models.CountryStats
		if err := rows.Scan(&cs.Country, &cs.CountryCode, &cs.Bans, &cs.Requests, &cs.Violations); err != nil {
			return nil, err
		}
		stats = append(stats, cs)
	}
	return stats, nil
}
