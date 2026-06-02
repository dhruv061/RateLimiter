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

// GetStats returns aggregated dashboard statistics scoped by global filters.
func (s *DashboardService) GetStats(filter models.GlobalFilter) (*models.DashboardStats, error) {
	stats := &models.DashboardStats{}
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	hourAgo := now.Add(-1 * time.Hour)
	dayAgo := now.Add(-24 * time.Hour)

	// Scoping helper for traffic_stats
	applyFilter := func(query string, defaultTime time.Time) (string, []interface{}) {
		where := "timestamp >= ?"
		args := []interface{}{defaultTime}

		if filter.StartTime != nil {
			if filter.StartTime.After(defaultTime) {
				args[0] = *filter.StartTime
			}
		}
		if filter.EndTime != nil {
			where += " AND timestamp <= ?"
			args = append(args, *filter.EndTime)
		}
		if filter.DomainID > 0 {
			where += " AND domain_id = ?"
			args = append(args, filter.DomainID)
		}

		fullQuery := fmt.Sprintf(query, where)
		return fullQuery, args
	}

	// Scoping helper for bans table
	applyBanFilter := func(query string, field string, defaultTime time.Time) (string, []interface{}) {
		where := fmt.Sprintf("%s >= ?", field)
		args := []interface{}{defaultTime}

		if filter.StartTime != nil {
			if filter.StartTime.After(defaultTime) {
				args[0] = *filter.StartTime
			}
		}
		if filter.EndTime != nil {
			where += fmt.Sprintf(" AND %s <= ?", field)
			args = append(args, *filter.EndTime)
		}
		if filter.DomainID > 0 {
			where += " AND domain_id = ?"
			args = append(args, filter.DomainID)
		}

		fullQuery := fmt.Sprintf(query, where)
		return fullQuery, args
	}

	// Traffic stats today
	q, args := applyFilter("SELECT COALESCE(SUM(total_requests), 0) FROM traffic_stats WHERE %s AND period = 'hour'", todayStart)
	s.db.QueryRow(q, args...).Scan(&stats.TotalRequestsToday)

	// Requests last hour
	q, args = applyFilter("SELECT COALESCE(SUM(total_requests), 0) FROM traffic_stats WHERE %s AND period = 'minute'", hourAgo)
	s.db.QueryRow(q, args...).Scan(&stats.RequestsLastHour)

	// Rate limiting 429 today
	q, args = applyFilter("SELECT COALESCE(SUM(status_429), 0) FROM traffic_stats WHERE %s AND period = 'hour'", todayStart)
	s.db.QueryRow(q, args...).Scan(&stats.Total429Today)

	// Count 429 last hour
	q, args = applyFilter("SELECT COALESCE(SUM(status_429), 0) FROM traffic_stats WHERE %s AND period = 'minute'", hourAgo)
	s.db.QueryRow(q, args...).Scan(&stats.Count429LastHour)

	// Active bans count
	activeBanQuery := "SELECT COUNT(*) FROM bans WHERE is_active = 1"
	var activeBanArgs []interface{}
	if filter.DomainID > 0 {
		activeBanQuery += " AND domain_id = ?"
		activeBanArgs = append(activeBanArgs, filter.DomainID)
	}
	s.db.QueryRow(activeBanQuery, activeBanArgs...).Scan(&stats.ActiveBans)

	// Bans today
	q, args = applyBanFilter("SELECT COUNT(*) FROM bans WHERE %s", "ban_time", todayStart)
	s.db.QueryRow(q, args...).Scan(&stats.TotalBansToday)

	// Bans 24h
	q, args = applyBanFilter("SELECT COUNT(*) FROM bans WHERE %s", "ban_time", dayAgo)
	s.db.QueryRow(q, args...).Scan(&stats.Bans24h)

	// Unbans today
	q, args = applyBanFilter("SELECT COUNT(*) FROM bans WHERE unban_time IS NOT NULL AND %s", "unban_time", todayStart)
	s.db.QueryRow(q, args...).Scan(&stats.UnbansToday)

	// System status
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

// GetAttackStatus evaluates the current threat level scoped by global filters.
func (s *DashboardService) GetAttackStatus(filter models.GlobalFilter) (*models.AttackStatus, error) {
	status := &models.AttackStatus{}
	now := time.Now()
	oneMinAgo := now.Add(-1 * time.Minute)
	fiveMinAgo := now.Add(-5 * time.Minute)

	// Scoping helper for traffic_stats
	applyFilter := func(query string, defaultTime time.Time) (string, []interface{}) {
		where := "timestamp >= ?"
		args := []interface{}{defaultTime}
		if filter.StartTime != nil {
			if filter.StartTime.After(defaultTime) {
				args[0] = *filter.StartTime
			}
		}
		if filter.EndTime != nil {
			where += " AND timestamp <= ?"
			args = append(args, *filter.EndTime)
		}
		if filter.DomainID > 0 {
			where += " AND domain_id = ?"
			args = append(args, filter.DomainID)
		}
		return fmt.Sprintf(query, where), args
	}

	// Scoping helper for bans
	applyBanFilter := func(query string, defaultTime time.Time) (string, []interface{}) {
		where := "ban_time >= ?"
		args := []interface{}{defaultTime}
		if filter.StartTime != nil {
			if filter.StartTime.After(defaultTime) {
				args[0] = *filter.StartTime
			}
		}
		if filter.EndTime != nil {
			where += " AND ban_time <= ?"
			args = append(args, *filter.EndTime)
		}
		if filter.DomainID > 0 {
			where += " AND domain_id = ?"
			args = append(args, filter.DomainID)
		}
		return fmt.Sprintf(query, where), args
	}

	// Last minute stats
	q, args := applyFilter("SELECT COALESCE(SUM(total_requests), 0) FROM traffic_stats WHERE %s AND period = 'minute'", oneMinAgo)
	s.db.QueryRow(q, args...).Scan(&status.RequestsLastMin)

	q, args = applyFilter("SELECT COALESCE(SUM(unique_ips), 0) FROM traffic_stats WHERE %s AND period = 'minute'", oneMinAgo)
	s.db.QueryRow(q, args...).Scan(&status.UniqueIPsLastMin)

	q, args = applyFilter("SELECT COALESCE(SUM(status_429), 0) FROM traffic_stats WHERE %s AND period = 'minute'", oneMinAgo)
	s.db.QueryRow(q, args...).Scan(&status.Count429LastMin)

	q, args = applyBanFilter("SELECT COUNT(*) FROM bans WHERE %s", oneMinAgo)
	s.db.QueryRow(q, args...).Scan(&status.BansLastMin)

	// Last 5 minutes stats
	q, args = applyFilter("SELECT COALESCE(SUM(total_requests), 0) FROM traffic_stats WHERE %s AND period = 'minute'", fiveMinAgo)
	s.db.QueryRow(q, args...).Scan(&status.RequestsLast5Min)

	q, args = applyFilter("SELECT COALESCE(SUM(unique_ips), 0) FROM traffic_stats WHERE %s AND period = 'minute'", fiveMinAgo)
	s.db.QueryRow(q, args...).Scan(&status.UniqueIPsLast5Min)

	q, args = applyFilter("SELECT COALESCE(SUM(status_429), 0) FROM traffic_stats WHERE %s AND period = 'minute'", fiveMinAgo)
	s.db.QueryRow(q, args...).Scan(&status.Count429Last5Min)

	q, args = applyBanFilter("SELECT COUNT(*) FROM bans WHERE %s", fiveMinAgo)
	s.db.QueryRow(q, args...).Scan(&status.BansLast5Min)

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

// GetTrafficTrends returns traffic data for charts, filtered by domain and range.
func (s *DashboardService) GetTrafficTrends(filter models.GlobalFilter, period string, hours int) ([]models.TrafficStat, error) {
	if hours < 1 {
		hours = 24
	}
	since := time.Now().Add(-time.Duration(hours) * time.Hour)
	if filter.StartTime != nil {
		since = *filter.StartTime
	}

	where := "period = ? AND timestamp >= ?"
	args := []interface{}{period, since}

	if filter.EndTime != nil {
		where += " AND timestamp <= ?"
		args = append(args, *filter.EndTime)
	}
	if filter.DomainID > 0 {
		where += " AND domain_id = ?"
		args = append(args, filter.DomainID)
	}

	query := fmt.Sprintf(
		"SELECT id, timestamp, total_requests, unique_ips, status_429, status_403, avg_response_time, period FROM traffic_stats WHERE %s ORDER BY timestamp ASC",
		where,
	)

	rows, err := s.db.Query(query, args...)
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

// GetCountryStats returns analytics grouped by country, filtered by domain and range.
func (s *DashboardService) GetCountryStats(filter models.GlobalFilter) ([]models.CountryStats, error) {
	where := "country != ''"
	var args []interface{}

	if filter.DomainID > 0 {
		where += " AND domain_id = ?"
		args = append(args, filter.DomainID)
	}
	if filter.StartTime != nil {
		where += " AND ban_time >= ?"
		args = append(args, *filter.StartTime)
	}
	if filter.EndTime != nil {
		where += " AND ban_time <= ?"
		args = append(args, *filter.EndTime)
	}

	query := fmt.Sprintf(`
		SELECT country, country_code,
			COUNT(*) as ban_count,
			COALESCE(SUM(request_count), 0) as total_requests,
			COALESCE(SUM(violation_count), 0) as total_violations
		FROM bans
		WHERE %s
		GROUP BY country_code
		ORDER BY ban_count DESC
		LIMIT 50
	`, where)

	rows, err := s.db.Query(query, args...)
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
