package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration.
type Config struct {
	// Server
	AppPort string
	AppEnv  string
	AppName string

	// Auth
	JWTSecret string
	JWTExpiry time.Duration

	// Database
	DatabasePath string

	// Log files
	NginxAccessLog string
	NginxErrorLog  string
	Fail2BanLog    string
	BlockFilePath  string

	// Demo mode
	DemoMode bool

	// Default admin
	DefaultAdminUser string
	DefaultAdminPass string

	// GeoIP
	GeoIPProvider string

	// Logging
	LogLevel string
	LogPath  string
}

// Load reads configuration from environment variables.
func Load() *Config {
	// Load .env file if it exists (ignore errors in production)
	_ = godotenv.Load()
	_ = godotenv.Load("../.env") // for running from backend/ dir

	expiry, err := time.ParseDuration(getEnv("JWT_EXPIRY", "24h"))
	if err != nil {
		expiry = 24 * time.Hour
	}

	return &Config{
		AppPort: getEnv("APP_PORT", "8080"),
		AppEnv:  getEnv("APP_ENV", "development"),
		AppName: getEnv("APP_NAME", "Fail2Ban Dashboard"),

		JWTSecret: getEnv("JWT_SECRET", "dev-secret-change-me"),
		JWTExpiry: expiry,

		DatabasePath: getEnv("DATABASE_PATH", "./data/dashboard.db"),

		NginxAccessLog: getEnv("NGINX_ACCESS_LOG", "/host/nginx/domain_access.log"),
		NginxErrorLog:  getEnv("NGINX_ERROR_LOG", "/host/nginx/domain_error.log"),
		Fail2BanLog:    getEnv("FAIL2BAN_LOG", "/host/fail2ban.log"),
		BlockFilePath:  getEnv("BLOCK_FILE_PATH", "/host/fail2ban_blocked.conf"),

		DemoMode: getEnv("DEMO_MODE", "false") == "true",

		DefaultAdminUser: getEnv("DEFAULT_ADMIN_USER", "admin"),
		DefaultAdminPass: getEnv("DEFAULT_ADMIN_PASS", "admin"),

		GeoIPProvider: getEnv("GEOIP_PROVIDER", "ip-api"),

		LogLevel: getEnv("LOG_LEVEL", "info"),
		LogPath:  getEnv("LOG_PATH", "./logs/dashboard.log"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
