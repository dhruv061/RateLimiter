package services

import (
	"database/sql"
	"fmt"
	"math"
	"time"

	"fail2ban-dashboard/internal/database"
	"fail2ban-dashboard/internal/models"
)

// BanService handles ban-related business logic.
type BanService struct {
	db *database.DB
}

// NewBanService creates a new BanService.
func NewBanService(db *database.DB) *BanService {
	return &BanService{db: db}
}

// GetActiveBans returns paginated active bans scoped by global filters.
func (s *BanService) GetActiveBans(filter models.GlobalFilter, page, perPage int, search, sortBy, sortDir, country string) (*models.PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	// Build query
	where := "WHERE is_active = 1"
	args := []interface{}{}

	if search != "" {
		where += " AND (ip_address LIKE ? OR country LIKE ? OR reason LIKE ?)"
		pattern := "%" + search + "%"
		args = append(args, pattern, pattern, pattern)
	}
	if country != "" {
		where += " AND country_code = ?"
		args = append(args, country)
	}

	where, args = applyGlobalFilter(where, args, filter)

	// Count
	var total int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM bans %s", where)
	if err := s.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, err
	}

	// Validate sort
	allowedSorts := map[string]bool{
		"ban_time": true, "ip_address": true, "country": true,
		"remaining_time": true, "request_count": true, "violation_count": true,
	}
	if !allowedSorts[sortBy] {
		sortBy = "ban_time"
	}
	if sortDir != "asc" {
		sortDir = "desc"
	}

	// For remaining_time, sort by ban_time + ban_duration
	orderBy := sortBy
	if sortBy == "remaining_time" {
		orderBy = "(ban_time + ban_duration)"
	}

	query := fmt.Sprintf(
		"SELECT id, ip_address, country, country_code, region, city, asn, isp, jail, reason, ban_time, unban_time, ban_duration, request_count, violation_count, is_active, domain_id, created_at FROM bans %s ORDER BY %s %s LIMIT ? OFFSET ?",
		where, orderBy, sortDir,
	)
	args = append(args, perPage, offset)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bans := []models.Ban{}
	for rows.Next() {
		var b models.Ban
		var unbanTime sql.NullTime
		err := rows.Scan(&b.ID, &b.IPAddress, &b.Country, &b.CountryCode, &b.Region, &b.City,
			&b.ASN, &b.ISP, &b.Jail, &b.Reason, &b.BanTime, &unbanTime,
			&b.BanDuration, &b.RequestCount, &b.ViolationCount, &b.IsActive, &b.DomainID, &b.CreatedAt)
		if err != nil {
			return nil, err
		}
		if unbanTime.Valid {
			b.UnbanTime = &unbanTime.Time
		}
		bans = append(bans, b)
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &models.PaginatedResponse{
		Data:       bans,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

// GetBanHistory returns paginated historical bans scoped by global filters.
func (s *BanService) GetBanHistory(filter models.GlobalFilter, page, perPage int, search, dateFrom, dateTo string) (*models.PaginatedResponse, error) {
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
		where += " AND (ip_address LIKE ? OR country LIKE ? OR reason LIKE ?)"
		pattern := "%" + search + "%"
		args = append(args, pattern, pattern, pattern)
	}
	if dateFrom != "" {
		where += " AND ban_time >= ?"
		args = append(args, dateFrom)
	} else if filter.StartTime != nil {
		where += " AND ban_time >= ?"
		args = append(args, *filter.StartTime)
	}

	if dateTo != "" {
		where += " AND ban_time <= ?"
		args = append(args, dateTo)
	} else if filter.EndTime != nil {
		where += " AND ban_time <= ?"
		args = append(args, *filter.EndTime)
	}

	if filter.DomainID > 0 {
		where += " AND domain_id = ?"
		args = append(args, filter.DomainID)
	}

	var total int64
	if err := s.db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM bans %s", where), args...).Scan(&total); err != nil {
		return nil, err
	}

	query := fmt.Sprintf(
		"SELECT id, ip_address, country, country_code, region, city, asn, isp, jail, reason, ban_time, unban_time, ban_duration, request_count, violation_count, is_active, domain_id, created_at FROM bans %s ORDER BY ban_time DESC LIMIT ? OFFSET ?",
		where,
	)
	args = append(args, perPage, offset)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bans := []models.Ban{}
	for rows.Next() {
		var b models.Ban
		var unbanTime sql.NullTime
		err := rows.Scan(&b.ID, &b.IPAddress, &b.Country, &b.CountryCode, &b.Region, &b.City,
			&b.ASN, &b.ISP, &b.Jail, &b.Reason, &b.BanTime, &unbanTime,
			&b.BanDuration, &b.RequestCount, &b.ViolationCount, &b.IsActive, &b.DomainID, &b.CreatedAt)
		if err != nil {
			return nil, err
		}
		if unbanTime.Valid {
			b.UnbanTime = &unbanTime.Time
		}
		bans = append(bans, b)
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &models.PaginatedResponse{
		Data:       bans,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

// GetBanByID returns a single ban by ID.
func (s *BanService) GetBanByID(id int64) (*models.Ban, error) {
	var b models.Ban
	var unbanTime sql.NullTime
	err := s.db.QueryRow(
		"SELECT id, ip_address, country, country_code, region, city, asn, isp, jail, reason, ban_time, unban_time, ban_duration, request_count, violation_count, is_active, domain_id, created_at FROM bans WHERE id = ?",
		id,
	).Scan(&b.ID, &b.IPAddress, &b.Country, &b.CountryCode, &b.Region, &b.City,
		&b.ASN, &b.ISP, &b.Jail, &b.Reason, &b.BanTime, &unbanTime,
		&b.BanDuration, &b.RequestCount, &b.ViolationCount, &b.IsActive, &b.DomainID, &b.CreatedAt)
	if err != nil {
		return nil, err
	}
	if unbanTime.Valid {
		b.UnbanTime = &unbanTime.Time
	}
	return &b, nil
}

// UnbanIP marks a ban as inactive.
func (s *BanService) UnbanIP(id int64) error {
	now := time.Now()
	_, err := s.db.Exec(
		"UPDATE bans SET is_active = 0, unban_time = ? WHERE id = ? AND is_active = 1",
		now, id,
	)
	return err
}

// BulkUnban marks multiple bans as inactive.
func (s *BanService) BulkUnban(ids []int64) (int64, error) {
	now := time.Now()
	var affected int64
	for _, id := range ids {
		result, err := s.db.Exec(
			"UPDATE bans SET is_active = 0, unban_time = ? WHERE id = ? AND is_active = 1",
			now, id,
		)
		if err != nil {
			return affected, err
		}
		n, _ := result.RowsAffected()
		affected += n
	}
	return affected, nil
}

// GetTopOffenders returns IPs with the most violations scoped by global filters.
func (s *BanService) GetTopOffenders(filter models.GlobalFilter, limit int) ([]models.TopOffender, error) {
	if limit < 1 || limit > 100 {
		limit = 10
	}

	where := "WHERE 1=1"
	var args []interface{}
	where, args = applyGlobalFilter(where, args, filter)

	query := fmt.Sprintf(`
		SELECT ip_address, country, country_code,
			SUM(request_count) as total_requests,
			SUM(violation_count) as total_violations,
			COUNT(*) as ban_count
		FROM bans
		%s
		GROUP BY ip_address
		ORDER BY total_violations DESC
		LIMIT ?
	`, where)
	args = append(args, limit)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var offenders []models.TopOffender
	for rows.Next() {
		var o models.TopOffender
		if err := rows.Scan(&o.IPAddress, &o.Country, &o.CountryCode, &o.TotalRequests, &o.ViolationCount, &o.BanCount); err != nil {
			return nil, err
		}
		offenders = append(offenders, o)
	}
	return offenders, nil
}

// GetActiveBanCount returns the count of active bans scoped by global filters.
func (s *BanService) GetActiveBanCount(filter models.GlobalFilter) (int, error) {
	where := "WHERE is_active = 1"
	var args []interface{}
	where, args = applyGlobalFilter(where, args, filter)

	var count int
	err := s.db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM bans %s", where), args...).Scan(&count)
	return count, err
}

// Helper function to append global filters (domain_id, start_time, end_time) to queries.
func applyGlobalFilter(where string, args []interface{}, filter models.GlobalFilter) (string, []interface{}) {
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
	return where, args
}
