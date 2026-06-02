package services

import (
	"fail2ban-dashboard/internal/database"
	"fail2ban-dashboard/internal/models"
	"fmt"
	"math"
	"time"
)

// WhitelistService handles whitelist operations.
type WhitelistService struct {
	db *database.DB
}

// NewWhitelistService creates a new WhitelistService.
func NewWhitelistService(db *database.DB) *WhitelistService {
	return &WhitelistService{db: db}
}

// GetAll returns all whitelist entries scoped by domain.
func (s *WhitelistService) GetAll(filter models.GlobalFilter, page, perPage int, search string) (*models.PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 50
	}
	offset := (page - 1) * perPage

	where := "WHERE 1=1"
	args := []interface{}{}
	if search != "" {
		where += " AND (ip_address LIKE ? OR description LIKE ?)"
		pattern := "%" + search + "%"
		args = append(args, pattern, pattern)
	}

	if filter.DomainID > 0 {
		where += " AND domain_id = ?"
		args = append(args, filter.DomainID)
	}

	var total int64
	s.db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM whitelist %s", where), args...).Scan(&total)

	query := fmt.Sprintf("SELECT id, ip_address, description, added_by, domain_id, created_at FROM whitelist %s ORDER BY created_at DESC LIMIT ? OFFSET ?", where)
	args = append(args, perPage, offset)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.WhitelistEntry
	for rows.Next() {
		var e models.WhitelistEntry
		if err := rows.Scan(&e.ID, &e.IPAddress, &e.Description, &e.AddedBy, &e.DomainID, &e.CreatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	return &models.PaginatedResponse{
		Data:       entries,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

// GetByID retrieves a single whitelist entry by its ID.
func (s *WhitelistService) GetByID(id int64) (*models.WhitelistEntry, error) {
	var e models.WhitelistEntry
	err := s.db.QueryRow("SELECT id, ip_address, description, added_by, domain_id, created_at FROM whitelist WHERE id = ?", id).Scan(
		&e.ID, &e.IPAddress, &e.Description, &e.AddedBy, &e.DomainID, &e.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// Add creates a new whitelist entry.
func (s *WhitelistService) Add(domainID int64, ip, description, addedBy string) (*models.WhitelistEntry, error) {
	result, err := s.db.Exec(
		"INSERT INTO whitelist (ip_address, description, added_by, domain_id) VALUES (?, ?, ?, ?)",
		ip, description, addedBy, domainID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add whitelist entry (IP may already exist): %w", err)
	}

	id, _ := result.LastInsertId()
	return &models.WhitelistEntry{
		ID:          id,
		IPAddress:   ip,
		Description: description,
		AddedBy:     addedBy,
		DomainID:    domainID,
		CreatedAt:   time.Now(),
	}, nil
}

// Remove deletes a whitelist entry by ID.
func (s *WhitelistService) Remove(id int64) error {
	result, err := s.db.Exec("DELETE FROM whitelist WHERE id = ?", id)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("whitelist entry not found")
	}
	return nil
}

// IsWhitelisted checks if an IP is in the whitelist for a specific domain.
func (s *WhitelistService) IsWhitelisted(domainID int64, ip string) bool {
	var count int
	s.db.QueryRow("SELECT COUNT(*) FROM whitelist WHERE ip_address = ? AND domain_id = ?", ip, domainID).Scan(&count)
	return count > 0
}
