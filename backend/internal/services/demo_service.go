package services

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"fail2ban-dashboard/internal/database"
)

// DemoService generates realistic fake data for demonstration.
type DemoService struct {
	db *database.DB
}

// NewDemoService creates a new DemoService.
func NewDemoService(db *database.DB) *DemoService {
	return &DemoService{db: db}
}

// countries with realistic distribution
var demoCountries = []struct {
	Name string
	Code string
	Weight int
}{
	{"China", "CN", 25},
	{"Russia", "RU", 20},
	{"United States", "US", 10},
	{"Brazil", "BR", 8},
	{"India", "IN", 7},
	{"Germany", "DE", 5},
	{"Netherlands", "NL", 5},
	{"France", "FR", 4},
	{"Vietnam", "VN", 4},
	{"Indonesia", "ID", 3},
	{"South Korea", "KR", 3},
	{"Ukraine", "UA", 3},
	{"United Kingdom", "GB", 2},
	{"Japan", "JP", 1},
}

var demoISPs = []string{
	"China Telecom", "Rostelecom", "DigitalOcean", "OVH SAS",
	"Amazon AWS", "Google Cloud", "Hetzner Online", "Linode",
	"Vultr Holdings", "Alibaba Cloud", "Tencent Cloud", "Azure",
}

var demoJails = []string{
	"nginx-429", "nginx-limit-req", "nginx-botsearch", "nginx-badbots",
}

var demoReasons = []string{
	"Rate limit exceeded (429)",
	"Too many requests in findtime",
	"Brute force attempt detected",
	"Bot scanning detected",
	"Excessive 429 responses",
	"DDoS pattern detected",
	"Automated scraping detected",
}

var demoUserAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
	"python-requests/2.31.0",
	"curl/8.4.0",
	"Go-http-client/2.0",
	"Googlebot/2.1 (+http://www.google.com/bot.html)",
	"Mozilla/5.0 (compatible; Bingbot/2.0)",
	"Scrapy/2.11",
	"axios/1.6.0",
}

var demoURLs = []string{
	"/api/v1/login",
	"/api/v1/users",
	"/api/v1/data",
	"/wp-admin",
	"/wp-login.php",
	"/.env",
	"/xmlrpc.php",
	"/api/v1/search",
	"/admin/config",
	"/phpmyadmin",
	"/api/v1/products",
	"/",
	"/index.html",
	"/robots.txt",
	"/sitemap.xml",
}

var demoAuditActions = []struct {
	Action string
	Target string
}{
	{"unban_ip", "192.168.1.%d"},
	{"add_whitelist", "10.0.0.%d"},
	{"remove_whitelist", "172.16.0.%d"},
	{"update_settings", "fail2ban_ban_time"},
	{"update_settings", "nginx_rate_limit_rps"},
	{"login", "admin"},
	{"reload_nginx", "nginx"},
	{"reload_fail2ban", "fail2ban"},
}

// SeedDemoData populates the database with realistic demo data.
func (s *DemoService) SeedDemoData() error {
	// Check if demo data already exists
	var count int
	s.db.QueryRow("SELECT COUNT(*) FROM bans").Scan(&count)
	if count > 10 {
		log.Println("Demo data already exists, skipping seed")
		return nil
	}

	log.Println("🌱 Seeding demo data...")

	if err := s.seedBans(); err != nil {
		return fmt.Errorf("failed to seed bans: %w", err)
	}
	if err := s.seedTrafficStats(); err != nil {
		return fmt.Errorf("failed to seed traffic stats: %w", err)
	}
	if err := s.seedWhitelist(); err != nil {
		return fmt.Errorf("failed to seed whitelist: %w", err)
	}
	if err := s.seedAuditLogs(); err != nil {
		return fmt.Errorf("failed to seed audit logs: %w", err)
	}

	log.Println("✅ Demo data seeded successfully")
	return nil
}

func (s *DemoService) seedBans() error {
	now := time.Now()
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO bans (ip_address, country, country_code, region, city, asn, isp, jail, reason,
			ban_time, unban_time, ban_duration, request_count, violation_count, is_active)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Generate ~45 active bans
	for i := 0; i < 45; i++ {
		c := weightedCountry()
		ip := randomIP()
		banTime := now.Add(-time.Duration(rand.Intn(7200)) * time.Second) // within last 2 hours
		duration := []int{3600, 7200, 14400, 86400}[rand.Intn(4)]

		_, err := stmt.Exec(
			ip, c.Name, c.Code, "Region-"+fmt.Sprint(rand.Intn(10)),
			"City-"+fmt.Sprint(rand.Intn(50)),
			fmt.Sprintf("AS%d", 10000+rand.Intn(90000)),
			demoISPs[rand.Intn(len(demoISPs))],
			demoJails[rand.Intn(len(demoJails))],
			demoReasons[rand.Intn(len(demoReasons))],
			banTime, nil, duration,
			50+rand.Intn(500), 5+rand.Intn(100), true,
		)
		if err != nil {
			return err
		}
	}

	// Generate ~200 historical bans (inactive)
	for i := 0; i < 200; i++ {
		c := weightedCountry()
		ip := randomIP()
		banTime := now.Add(-time.Duration(rand.Intn(720)) * time.Hour) // within last 30 days
		duration := []int{3600, 7200, 14400, 86400}[rand.Intn(4)]
		unbanTime := banTime.Add(time.Duration(duration) * time.Second)

		_, err := stmt.Exec(
			ip, c.Name, c.Code, "Region-"+fmt.Sprint(rand.Intn(10)),
			"City-"+fmt.Sprint(rand.Intn(50)),
			fmt.Sprintf("AS%d", 10000+rand.Intn(90000)),
			demoISPs[rand.Intn(len(demoISPs))],
			demoJails[rand.Intn(len(demoJails))],
			demoReasons[rand.Intn(len(demoReasons))],
			banTime, unbanTime, duration,
			50+rand.Intn(500), 5+rand.Intn(100), false,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *DemoService) seedTrafficStats() error {
	now := time.Now()
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO traffic_stats (timestamp, total_requests, unique_ips, status_429, status_403, avg_response_time, period)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Minute-level stats for last 2 hours
	for i := 0; i < 120; i++ {
		ts := now.Add(-time.Duration(i) * time.Minute)
		baseReqs := 100 + rand.Intn(400)
		// Simulate traffic spike
		if i >= 30 && i <= 45 {
			baseReqs = 500 + rand.Intn(1000)
		}
		stmt.Exec(ts, baseReqs, 20+rand.Intn(80), rand.Intn(baseReqs/10), rand.Intn(5), 50+rand.Float64()*200, "minute")
	}

	// Hour-level stats for last 7 days
	for i := 0; i < 168; i++ {
		ts := now.Add(-time.Duration(i) * time.Hour)
		baseReqs := 5000 + rand.Intn(15000)
		// Day/night pattern
		hour := ts.Hour()
		if hour >= 2 && hour <= 6 {
			baseReqs = 2000 + rand.Intn(3000) // Low traffic at night
		}
		if hour >= 10 && hour <= 16 {
			baseReqs = 10000 + rand.Intn(10000) // High traffic during day
		}
		stmt.Exec(ts, baseReqs, 500+rand.Intn(2000), rand.Intn(baseReqs/20), rand.Intn(50), 30+rand.Float64()*150, "hour")
	}

	// Day-level stats for last 30 days
	for i := 0; i < 30; i++ {
		ts := now.Add(-time.Duration(i) * 24 * time.Hour)
		baseReqs := 100000 + rand.Intn(200000)
		stmt.Exec(ts, baseReqs, 5000+rand.Intn(10000), rand.Intn(baseReqs/50), rand.Intn(200), 40+rand.Float64()*120, "day")
	}

	return tx.Commit()
}

func (s *DemoService) seedWhitelist() error {
	entries := []struct{ IP, Desc string }{
		{"10.0.0.1", "Internal monitoring server"},
		{"192.168.1.1", "Office gateway"},
		{"203.0.113.50", "Partner API server"},
		{"198.51.100.10", "CDN origin"},
		{"172.16.0.100", "Staging server"},
		{"8.8.8.8", "Google DNS (testing)"},
	}

	for _, e := range entries {
		s.db.Exec(
			"INSERT OR IGNORE INTO whitelist (ip_address, description, added_by) VALUES (?, ?, 'admin')",
			e.IP, e.Desc,
		)
	}
	return nil
}

func (s *DemoService) seedAuditLogs() error {
	now := time.Now()
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(
		"INSERT INTO audit_logs (user_id, username, action, target, details, ip_address, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i := 0; i < 50; i++ {
		ts := now.Add(-time.Duration(rand.Intn(720)) * time.Hour)
		a := demoAuditActions[rand.Intn(len(demoAuditActions))]
		target := fmt.Sprintf(a.Target, rand.Intn(255))
		stmt.Exec(1, "admin", a.Action, target, "{}", "127.0.0.1", ts)
	}

	return tx.Commit()
}

func weightedCountry() struct{ Name, Code string } {
	totalWeight := 0
	for _, c := range demoCountries {
		totalWeight += c.Weight
	}
	r := rand.Intn(totalWeight)
	for _, c := range demoCountries {
		r -= c.Weight
		if r < 0 {
			return struct{ Name, Code string }{c.Name, c.Code}
		}
	}
	return struct{ Name, Code string }{demoCountries[0].Name, demoCountries[0].Code}
}

func randomIP() string {
	return fmt.Sprintf("%d.%d.%d.%d", 1+rand.Intn(223), rand.Intn(256), rand.Intn(256), 1+rand.Intn(254))
}
