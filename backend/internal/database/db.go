package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// DB wraps the sql.DB connection.
type DB struct {
	*sql.DB
}

// New creates a new database connection and runs migrations.
func New(dbPath string) (*DB, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Connection pool settings for SQLite
	db.SetMaxOpenConns(1) // SQLite doesn't support concurrent writes
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	wrapped := &DB{db}

	if err := wrapped.migrate(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("✅ Database initialized successfully")
	return wrapped, nil
}

// migrate runs all database migrations.
func (db *DB) migrate() error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			email TEXT DEFAULT '',
			role TEXT DEFAULT 'admin',
			must_change_pass BOOLEAN DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			last_login DATETIME
		)`,

		`CREATE TABLE IF NOT EXISTS bans (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ip_address TEXT NOT NULL,
			country TEXT DEFAULT '',
			country_code TEXT DEFAULT '',
			region TEXT DEFAULT '',
			city TEXT DEFAULT '',
			asn TEXT DEFAULT '',
			isp TEXT DEFAULT '',
			jail TEXT DEFAULT '',
			reason TEXT DEFAULT '',
			ban_time DATETIME NOT NULL,
			unban_time DATETIME,
			ban_duration INTEGER DEFAULT 3600,
			request_count INTEGER DEFAULT 0,
			violation_count INTEGER DEFAULT 0,
			is_active BOOLEAN DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS whitelist (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ip_address TEXT UNIQUE NOT NULL,
			description TEXT DEFAULT '',
			added_by TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS audit_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			username TEXT DEFAULT '',
			action TEXT NOT NULL,
			target TEXT DEFAULT '',
			details TEXT DEFAULT '',
			ip_address TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,

		`CREATE TABLE IF NOT EXISTS settings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			key TEXT UNIQUE NOT NULL,
			value TEXT DEFAULT '',
			updated_by TEXT DEFAULT '',
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS traffic_stats (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp DATETIME NOT NULL,
			total_requests INTEGER DEFAULT 0,
			unique_ips INTEGER DEFAULT 0,
			status_429 INTEGER DEFAULT 0,
			status_403 INTEGER DEFAULT 0,
			avg_response_time REAL DEFAULT 0,
			period TEXT NOT NULL
		)`,

		// Indexes
		`CREATE INDEX IF NOT EXISTS idx_bans_ip ON bans(ip_address)`,
		`CREATE INDEX IF NOT EXISTS idx_bans_active ON bans(is_active)`,
		`CREATE INDEX IF NOT EXISTS idx_bans_ban_time ON bans(ban_time)`,
		`CREATE INDEX IF NOT EXISTS idx_bans_active_time ON bans(is_active, ban_time)`,
		`CREATE INDEX IF NOT EXISTS idx_audit_created ON audit_logs(created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_traffic_timestamp ON traffic_stats(timestamp, period)`,
		`CREATE INDEX IF NOT EXISTS idx_whitelist_ip ON whitelist(ip_address)`,
	}

	for _, m := range migrations {
		if _, err := db.Exec(m); err != nil {
			return fmt.Errorf("migration failed: %w\nSQL: %s", err, m)
		}
	}

	return nil
}

// SeedDefaultAdmin creates the default admin user if no users exist.
func (db *DB) SeedDefaultAdmin(username, password string) error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil // Users already exist
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	_, err = db.Exec(
		`INSERT INTO users (username, password_hash, role, must_change_pass) VALUES (?, ?, 'admin', 1)`,
		username, string(hash),
	)
	if err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	log.Printf("✅ Default admin user '%s' created (password change required on first login)", username)
	return nil
}

// SeedDefaultSettings inserts default settings if none exist.
func (db *DB) SeedDefaultSettings() error {
	defaults := map[string]string{
		"nginx_rate_limit_rps":     "10",
		"nginx_rate_limit_burst":   "20",
		"fail2ban_ban_time":        "3600",
		"fail2ban_find_time":       "600",
		"fail2ban_max_retry":       "5",
		"theme":                    "dark",
		"notifications_enabled":    "true",
		"auto_refresh_interval":    "30",
	}

	for key, value := range defaults {
		_, err := db.Exec(
			`INSERT OR IGNORE INTO settings (key, value, updated_by) VALUES (?, ?, 'system')`,
			key, value,
		)
		if err != nil {
			return fmt.Errorf("failed to seed setting %s: %w", key, err)
		}
	}

	return nil
}
