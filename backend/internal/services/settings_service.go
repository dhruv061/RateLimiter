package services

import (
	"fail2ban-dashboard/internal/database"
	"fail2ban-dashboard/internal/models"
)

// SettingsService handles application settings.
type SettingsService struct {
	db *database.DB
}

// NewSettingsService creates a new SettingsService.
func NewSettingsService(db *database.DB) *SettingsService {
	return &SettingsService{db: db}
}

// GetAll returns all settings as a map.
func (s *SettingsService) GetAll() (map[string]string, error) {
	rows, err := s.db.Query("SELECT key, value FROM settings")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	settings := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		settings[key] = value
	}
	return settings, nil
}

// Get returns a single setting value.
func (s *SettingsService) Get(key string) (string, error) {
	var value string
	err := s.db.QueryRow("SELECT value FROM settings WHERE key = ?", key).Scan(&value)
	return value, err
}

// Update updates multiple settings at once.
func (s *SettingsService) Update(settings map[string]string, updatedBy string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT OR REPLACE INTO settings (key, value, updated_by, updated_at) VALUES (?, ?, ?, CURRENT_TIMESTAMP)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for key, value := range settings {
		if _, err := stmt.Exec(key, value, updatedBy); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetSettingsDetailed returns settings with metadata.
func (s *SettingsService) GetSettingsDetailed() ([]models.Setting, error) {
	rows, err := s.db.Query("SELECT id, key, value, updated_by, updated_at FROM settings ORDER BY key")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []models.Setting
	for rows.Next() {
		var st models.Setting
		if err := rows.Scan(&st.ID, &st.Key, &st.Value, &st.UpdatedBy, &st.UpdatedAt); err != nil {
			return nil, err
		}
		settings = append(settings, st)
	}
	return settings, nil
}
