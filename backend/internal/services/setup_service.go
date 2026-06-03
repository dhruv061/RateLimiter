package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"fail2ban-dashboard/internal/database"
	"fail2ban-dashboard/internal/models"
)

// SetupService handles setup wizard operations.
type SetupService struct {
	db       *database.DB
	demoMode bool
}

// NewSetupService creates a new SetupService instance.
func NewSetupService(db *database.DB, demoMode bool) *SetupService {
	return &SetupService{db: db, demoMode: demoMode}
}

// domainToSlug converts example.com -> example-com
func domainToSlug(domain string) string {
	slug := strings.ReplaceAll(domain, ".", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	var clean []rune
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			clean = append(clean, r)
		} else if r >= 'A' && r <= 'Z' {
			clean = append(clean, r+32)
		}
	}
	return string(clean)
}

// domainToVarName converts domain slug to a valid Nginx variable name
func domainToVarName(slug string) string {
	return strings.ReplaceAll(slug, "-", "_")
}

// fileExists checks if a file exists on the host or inside container
func fileExists(path string) bool {
	hostPath := path
	if strings.HasPrefix(path, "/etc/fail2ban/") {
		hostPath = strings.Replace(path, "/etc/fail2ban/", "/host/fail2ban/", 1)
	} else if strings.HasPrefix(path, "/etc/nginx/") {
		hostPath = strings.Replace(path, "/etc/nginx/", "/host/nginx-config/", 1)
		if _, err := os.Stat(hostPath); os.IsNotExist(err) {
			if _, err := os.Stat(path); err == nil {
				return true
			}
			return false
		}
	}

	if _, err := os.Stat(hostPath); err == nil {
		return true
	}
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

// GetFail2BanStatus returns the current status of the Fail2Ban service.
func (s *SetupService) GetFail2BanStatus() (*models.Fail2BanStatusResponse, error) {
	if s.demoMode {
		return &models.Fail2BanStatusResponse{
			Installed:   true,
			Running:     true,
			Version:     "0.11.2",
			ActiveJails: []string{"nginx-limit-req", "nginx-429"},
			JailCount:   2,
		}, nil
	}

	resp := &models.Fail2BanStatusResponse{
		Installed:   false,
		Running:     false,
		Version:     "",
		ActiveJails: []string{},
		JailCount:   0,
	}

	_, err := exec.LookPath("fail2ban-client")
	hasClient := err == nil

	hasEtcConfig := false
	if _, err := os.Stat("/host/fail2ban"); err == nil {
		hasEtcConfig = true
	} else if _, err := os.Stat("/etc/fail2ban"); err == nil {
		hasEtcConfig = true
	}

	if !hasClient && !hasEtcConfig {
		return resp, nil
	}

	resp.Installed = true

	if hasClient {
		cmd := exec.Command("fail2ban-client", "status")
		output, err := cmd.CombinedOutput()
		if err == nil {
			resp.Running = true
			outStr := string(output)

			reCount := regexp.MustCompile(`(?i)Number of jail:\s*(\d+)`)
			if matches := reCount.FindStringSubmatch(outStr); len(matches) > 1 {
				if count, err := strconv.Atoi(matches[1]); err == nil {
					resp.JailCount = count
				}
			}

			reJails := regexp.MustCompile(`(?i)Jail list:\s*([^\r\n]+)`)
			if matches := reJails.FindStringSubmatch(outStr); len(matches) > 1 {
				jailsRaw := strings.Split(matches[1], ",")
				for _, j := range jailsRaw {
					resp.ActiveJails = append(resp.ActiveJails, strings.TrimSpace(j))
				}
			}
		}

		vCmd := exec.Command("fail2ban-client", "--version")
		vOutput, err := vCmd.CombinedOutput()
		if err == nil {
			vStr := string(vOutput)
			reVer := regexp.MustCompile(`(?i)v(\d+\.\d+\.\d+)`)
			if matches := reVer.FindStringSubmatch(vStr); len(matches) > 1 {
				resp.Version = matches[1]
			} else {
				lines := strings.Split(vStr, "\n")
				if len(lines) > 0 {
					resp.Version = strings.TrimSpace(lines[0])
				}
			}
		}
	}

	if !resp.Running && hasEtcConfig {
		logPath := "/host/fail2ban.log"
		if _, err := os.Stat(logPath); err != nil {
			logPath = "/var/log/fail2ban.log"
		}
		if _, err := os.Stat(logPath); err == nil {
			resp.Running = true
		}

		resp.Version = "0.11.x (CLI inaccessible)"

		jailDPath := "/host/fail2ban/jail.d"
		if _, err := os.Stat(jailDPath); err != nil {
			jailDPath = "/etc/fail2ban/jail.d"
		}

		if files, err := os.ReadDir(jailDPath); err == nil {
			for _, f := range files {
				if !f.IsDir() && (strings.HasSuffix(f.Name(), ".local") || strings.HasSuffix(f.Name(), ".conf")) {
					content, err := os.ReadFile(filepath.Join(jailDPath, f.Name()))
					if err == nil {
						reHeader := regexp.MustCompile(`(?m)^\[([^\]]+)\]`)
						matches := reHeader.FindAllStringSubmatch(string(content), -1)
						for _, m := range matches {
							if len(m) > 1 && !strings.Contains(m[1], "DEFAULT") {
								jailName := strings.TrimSpace(m[1])
								exists := false
								for _, existing := range resp.ActiveJails {
									if existing == jailName {
										exists = true
										break
									}
								}
								if !exists {
									resp.ActiveJails = append(resp.ActiveJails, jailName)
								}
							}
						}
					}
				}
			}
		}
		resp.JailCount = len(resp.ActiveJails)
	}

	return resp, nil
}

// DiscoverNginxDomains finds potential domains inside active Nginx sites.
func (s *SetupService) DiscoverNginxDomains() (*models.NginxDomainDiscoveryResponse, error) {
	if s.demoMode {
		return &models.NginxDomainDiscoveryResponse{
			Domains: []models.NginxDiscoveredDomain{
				{ServerName: "example.com", ConfigFile: "/etc/nginx/sites-available/example.com.conf", HasSSL: true},
				{ServerName: "api.example.com", ConfigFile: "/etc/nginx/sites-available/api.example.com.conf", HasSSL: true},
				{ServerName: "test.example.com", ConfigFile: "/etc/nginx/sites-available/test.example.com.conf", HasSSL: false},
			},
			ScannedDirs: []string{"/etc/nginx/sites-enabled", "/etc/nginx/conf.d", "/etc/nginx/sites-available"},
		}, nil
	}

	dirsToScan := []string{}
	hostDirs := []string{
		"/host/nginx-sites-enabled",
		"/host/nginx-conf-d",
		"/host/nginx-sites-available",
	}
	stdDirs := []string{
		"/etc/nginx/sites-enabled",
		"/etc/nginx/conf.d",
		"/etc/nginx/sites-available",
	}

	for i, hostDir := range hostDirs {
		if _, err := os.Stat(hostDir); err == nil {
			dirsToScan = append(dirsToScan, hostDir)
		} else {
			if _, err := os.Stat(stdDirs[i]); err == nil {
				dirsToScan = append(dirsToScan, stdDirs[i])
			}
		}
	}

	if len(dirsToScan) == 0 {
		dirsToScan = stdDirs
	}

	discovered := []models.NginxDiscoveredDomain{}
	seenDomains := map[string]bool{}

	serverNameRegex := regexp.MustCompile(`(?i)\bserver_name\s+([^;]+);`)
	sslRegex := regexp.MustCompile(`(?i)\b(listen\s+\d+\s+ssl|ssl_certificate\b)`)

	for _, dir := range dirsToScan {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, f := range files {
			if f.IsDir() {
				continue
			}
			filePath := filepath.Join(dir, f.Name())
			contentBytes, err := os.ReadFile(filePath)
			if err != nil {
				continue
			}
			content := string(contentBytes)

			matches := serverNameRegex.FindAllStringSubmatch(content, -1)
			hasSSL := sslRegex.MatchString(content)

			for _, match := range matches {
				if len(match) < 2 {
					continue
				}
				domains := strings.Fields(match[1])
				for _, d := range domains {
					d = strings.Trim(d, `;"'`)
					if d == "" || d == "_" || d == "localhost" || strings.HasPrefix(d, "$") {
						continue
					}
					if seenDomains[d] {
						continue
					}
					seenDomains[d] = true

					displayFile := filePath
					if strings.HasPrefix(filePath, "/host/") {
						displayFile = "/etc/" + strings.TrimPrefix(filePath, "/host/nginx-")
					}

					discovered = append(discovered, models.NginxDiscoveredDomain{
						ServerName: d,
						ConfigFile: displayFile,
						HasSSL:     hasSSL,
					})
				}
			}
		}
	}

	return &models.NginxDomainDiscoveryResponse{
		Domains:     discovered,
		ScannedDirs: dirsToScan,
	}, nil
}

// GenerateConfig generates the full set of Fail2Ban configurations and a setup script.
func (s *SetupService) GenerateConfig(domainName string, rateLimit, burst, banTime int) (*models.GeneratedConfigResponse, error) {
	if rateLimit <= 0 {
		rateLimit = 5
	}
	if burst <= 0 {
		burst = 5
	}
	if banTime <= 0 {
		banTime = 86400
	}

	slug := domainToSlug(domainName)
	varName := domainToVarName(slug)

	filterContent := `[Definition]
failregex = ^<HOST> -.* 429 
ignoreregex =
`

	actionContent := fmt.Sprintf(`[Definition]
actionban = echo "<ip> 1;" >> /etc/nginx/%s_blocked.conf && nginx -t && nginx -s reload
actionunban = sed -i '/<ip> 1;/d' /etc/nginx/%s_blocked.conf && nginx -t && nginx -s reload
`, slug, slug)

	jailContent := fmt.Sprintf(`[%s-limit-req]
enabled  = true
port     = http,https
filter   = nginx-limit-req
logpath  = /var/log/nginx/%s_error.log
backend  = auto
action   = %s-block
maxretry = 1
findtime = 5
bantime  = %d

[%s-429]
enabled  = true
port     = http,https
filter   = %s-429
logpath  = /var/log/nginx/%s_access.log
backend  = auto
action   = %s-block
maxretry = 1
findtime = 5
bantime  = %d
`, slug, slug, slug, banTime, slug, slug, slug, slug, banTime)

	blockContent := fmt.Sprintf(`# Blocked IPs for %s
# Fail2ban will automatically write "deny <ip>;" lines into this file.
# Do not edit this file manually.
`, domainName)

	nginxSnippet := fmt.Sprintf(`# ── Fail2ban auto-blocked IPs check ──
if ($%s_blocked) {
    return 403;
}

# ── Rate Limiting — %d req/sec, burst of %d ──
limit_req zone=%s_rate_limit burst=%d nodelay;
limit_req_status 429;
`, varName, rateLimit, burst, slug, burst)

	nginxZoneLine := fmt.Sprintf(`limit_req_zone $binary_remote_addr zone=%s_rate_limit:10m rate=%dr/s;

geo $%s_blocked {
    default 0;
    include /etc/nginx/%s_blocked.conf;
}
`, slug, rateLimit, varName, slug)

	setupScript := fmt.Sprintf(`#!/bin/bash
# Auto-generated setup script for %s by ShieldWatch Dashboard.
set -e

echo "======================================================================"
echo "Installing/configuring Fail2Ban & Nginx protection for %s"
echo "======================================================================"

# Create target directories if they do not exist
sudo mkdir -p /etc/fail2ban/filter.d /etc/fail2ban/action.d /etc/fail2ban/jail.d /etc/nginx

# 1. Create block file
echo "1. Creating empty block file at /etc/nginx/%s_blocked.conf"
sudo touch /etc/nginx/%s_blocked.conf
sudo chmod 644 /etc/nginx/%s_blocked.conf

# 2. Create filter file
echo "2. Writing Fail2Ban filter to /etc/fail2ban/filter.d/%s-429.conf"
sudo tee /etc/fail2ban/filter.d/%s-429.conf > /dev/null << 'EOF'
[Definition]
failregex = ^<HOST> -.* 429 
ignoreregex =
EOF

# 3. Create action file
echo "3. Writing Fail2Ban action to /etc/fail2ban/action.d/%s-block.conf"
sudo tee /etc/fail2ban/action.d/%s-block.conf > /dev/null << 'EOF'
[Definition]
actionban = echo "<ip> 1;" >> /etc/nginx/%s_blocked.conf && nginx -t && nginx -s reload
actionunban = sed -i '/<ip> 1;/d' /etc/nginx/%s_blocked.conf && nginx -t && nginx -s reload
EOF

# 4. Create jail file
echo "4. Writing Fail2Ban jail to /etc/fail2ban/jail.d/%s.local"
sudo tee /etc/fail2ban/jail.d/%s.local > /dev/null << 'EOF'
[%s-limit-req]
enabled  = true
port     = http,https
filter   = nginx-limit-req
logpath  = /var/log/nginx/%s_error.log
backend  = auto
action   = %s-block
maxretry = 1
findtime = 5
bantime  = %d

[%s-429]
enabled  = true
port     = http,https
filter   = %s-429
logpath  = /var/log/nginx/%s_access.log
backend  = auto
action   = %s-block
maxretry = 1
findtime = 5
bantime  = %d
EOF

# 5. Restart Services
echo "5. Restarting Fail2Ban service..."
sudo systemctl restart fail2ban || sudo service fail2ban restart

echo "6. Validating Nginx configuration..."
sudo nginx -t

echo "7. Reloading Nginx service..."
sudo nginx -s reload || sudo systemctl reload nginx

echo "======================================================================"
echo "Configuration files written & services reloaded successfully!"
echo "Please ensure you paste the generated Nginx snippet into your server block."
echo "======================================================================"
`, domainName, domainName, slug, slug, slug, slug, slug, slug, slug, slug, slug, slug, slug, slug, slug, slug, banTime, slug, slug, slug, slug, banTime)

	files := []models.GeneratedFile{
		{
			Filename: fmt.Sprintf("%s-429.conf", slug),
			Path:     fmt.Sprintf("/etc/fail2ban/filter.d/%s-429.conf", slug),
			Content:  filterContent,
			Type:     "filter",
		},
		{
			Filename: fmt.Sprintf("%s-block.conf", slug),
			Path:     fmt.Sprintf("/etc/fail2ban/action.d/%s-block.conf", slug),
			Content:  actionContent,
			Type:     "action",
		},
		{
			Filename: fmt.Sprintf("%s.local", slug),
			Path:     fmt.Sprintf("/etc/fail2ban/jail.d/%s.local", slug),
			Content:  jailContent,
			Type:     "jail",
		},
		{
			Filename: fmt.Sprintf("%s_blocked.conf", slug),
			Path:     fmt.Sprintf("/etc/nginx/%s_blocked.conf", slug),
			Content:  blockContent,
			Type:     "block",
		},
	}

	return &models.GeneratedConfigResponse{
		DomainSlug:    slug,
		Files:         files,
		SetupScript:   setupScript,
		NginxSnippet:  nginxSnippet,
		NginxZoneLine: nginxZoneLine,
	}, nil
}

// GenerateCleanupScript generates a removal script.
func (s *SetupService) GenerateCleanupScript(domainName string) (*models.CleanupScriptResponse, error) {
	slug := domainToSlug(domainName)

	cleanupScript := fmt.Sprintf(`#!/bin/bash
# Auto-generated cleanup script for %s by ShieldWatch Dashboard.
set -e

echo "======================================================================"
echo "Removing Fail2Ban & Nginx configurations for %s"
echo "======================================================================"

# 1. Remove Fail2Ban jail
if [ -f /etc/fail2ban/jail.d/%s.local ]; then
    echo "1. Removing Fail2Ban jail /etc/fail2ban/jail.d/%s.local"
    sudo rm -f /etc/fail2ban/jail.d/%s.local
fi

# 2. Remove filter
if [ -f /etc/fail2ban/filter.d/%s-429.conf ]; then
    echo "2. Removing Fail2Ban filter /etc/fail2ban/filter.d/%s-429.conf"
    sudo rm -f /etc/fail2ban/filter.d/%s-429.conf
fi

# 3. Remove action
if [ -f /etc/fail2ban/action.d/%s-block.conf ]; then
    echo "3. Removing Fail2Ban action /etc/fail2ban/action.d/%s-block.conf"
    sudo rm -f /etc/fail2ban/action.d/%s-block.conf
fi

# 4. Remove block file
if [ -f /etc/nginx/%s_blocked.conf ]; then
    echo "4. Removing Nginx block file /etc/nginx/%s_blocked.conf"
    sudo rm -f /etc/nginx/%s_blocked.conf
fi

# 5. Restart/Reload services
echo "5. Restarting Fail2Ban..."
sudo systemctl restart fail2ban || sudo service fail2ban restart

echo "6. Testing Nginx configuration..."
sudo nginx -t

echo "7. Reloading Nginx..."
sudo nginx -s reload || sudo systemctl reload nginx

echo "======================================================================"
echo "Cleanup script executed successfully!"
echo "Make sure you have manually removed the Nginx snippets from your config!"
echo "======================================================================"
`, domainName, domainName, slug, slug, slug, slug, slug, slug, slug, slug, slug, slug, slug)

	filesToRemove := []models.GeneratedFile{
		{Filename: fmt.Sprintf("%s.local", slug), Path: fmt.Sprintf("/etc/fail2ban/jail.d/%s.local", slug)},
		{Filename: fmt.Sprintf("%s-429.conf", slug), Path: fmt.Sprintf("/etc/fail2ban/filter.d/%s-429.conf", slug)},
		{Filename: fmt.Sprintf("%s-block.conf", slug), Path: fmt.Sprintf("/etc/fail2ban/action.d/%s-block.conf", slug)},
		{Filename: fmt.Sprintf("%s_blocked.conf", slug), Path: fmt.Sprintf("/etc/nginx/%s_blocked.conf", slug)},
	}

	return &models.CleanupScriptResponse{
		DomainSlug:    slug,
		CleanupScript: cleanupScript,
		FilesToRemove: filesToRemove,
	}, nil
}

// ValidateSetup validates that setup configs exist and are working.
func (s *SetupService) ValidateSetup(domainName string) (*models.SetupValidationResponse, error) {
	if s.demoMode {
		return &models.SetupValidationResponse{
			Checks: []models.SetupValidationCheck{
				{Name: "Fail2Ban Filter File", Passed: true, Message: "Filter /etc/fail2ban/filter.d/" + domainToSlug(domainName) + "-429.conf exists"},
				{Name: "Fail2Ban Action File", Passed: true, Message: "Action /etc/fail2ban/action.d/" + domainToSlug(domainName) + "-block.conf exists"},
				{Name: "Fail2Ban Jail File", Passed: true, Message: "Jail /etc/fail2ban/jail.d/" + domainToSlug(domainName) + ".local exists"},
				{Name: "Nginx Block File", Passed: true, Message: "Block file /etc/nginx/" + domainToSlug(domainName) + "_blocked.conf exists"},
				{Name: "Fail2Ban Jails Active", Passed: true, Message: "Jails " + domainToSlug(domainName) + "-limit-req and " + domainToSlug(domainName) + "-429 are running"},
			},
			OverallValid: true,
		}, nil
	}

	slug := domainToSlug(domainName)
	checks := []models.SetupValidationCheck{}
	overallValid := true

	filterPath := fmt.Sprintf("/etc/fail2ban/filter.d/%s-429.conf", slug)
	filterPassed := fileExists(filterPath)
	filterMsg := fmt.Sprintf("Filter file %s exists", filterPath)
	if !filterPassed {
		filterMsg = fmt.Sprintf("Filter file %s not found. Please run the setup script.", filterPath)
		overallValid = false
	}
	checks = append(checks, models.SetupValidationCheck{Name: "Fail2Ban Filter File", Passed: filterPassed, Message: filterMsg})

	actionPath := fmt.Sprintf("/etc/fail2ban/action.d/%s-block.conf", slug)
	actionPassed := fileExists(actionPath)
	actionMsg := fmt.Sprintf("Action file %s exists", actionPath)
	if !actionPassed {
		actionMsg = fmt.Sprintf("Action file %s not found. Please run the setup script.", actionPath)
		overallValid = false
	}
	checks = append(checks, models.SetupValidationCheck{Name: "Fail2Ban Action File", Passed: actionPassed, Message: actionMsg})

	jailPath := fmt.Sprintf("/etc/fail2ban/jail.d/%s.local", slug)
	jailPassed := fileExists(jailPath)
	jailMsg := fmt.Sprintf("Jail file %s exists", jailPath)
	if !jailPassed {
		jailMsg = fmt.Sprintf("Jail file %s not found. Please run the setup script.", jailPath)
		overallValid = false
	}
	checks = append(checks, models.SetupValidationCheck{Name: "Fail2Ban Jail File", Passed: jailPassed, Message: jailMsg})

	blockPath := fmt.Sprintf("/etc/nginx/%s_blocked.conf", slug)
	blockPassed := fileExists(blockPath)
	blockMsg := fmt.Sprintf("Nginx block file %s exists", blockPath)
	if !blockPassed {
		blockMsg = fmt.Sprintf("Block file %s not found. Please run the setup script.", blockPath)
		overallValid = false
	}
	checks = append(checks, models.SetupValidationCheck{Name: "Nginx Block File", Passed: blockPassed, Message: blockMsg})

	jailsActive := false
	jailActiveMsg := ""

	status, err := s.GetFail2BanStatus()
	if err == nil && status.Running {
		limitReqActive := false
		req429Active := false
		for _, activeJail := range status.ActiveJails {
			if activeJail == fmt.Sprintf("%s-limit-req", slug) {
				limitReqActive = true
			}
			if activeJail == fmt.Sprintf("%s-429", slug) {
				req429Active = true
			}
		}
		if limitReqActive && req429Active {
			jailsActive = true
			jailActiveMsg = fmt.Sprintf("Jails %s-limit-req and %s-429 are running", slug, slug)
		} else if limitReqActive {
			jailActiveMsg = fmt.Sprintf("Jail %s-limit-req is running, but %s-429 is not active.", slug, slug)
			overallValid = false
		} else if req429Active {
			jailActiveMsg = fmt.Sprintf("Jail %s-429 is running, but %s-limit-req is not active.", slug, slug)
			overallValid = false
		} else {
			jailActiveMsg = fmt.Sprintf("Neither %s-limit-req nor %s-429 jails are active in Fail2Ban. Ensure Fail2Ban restarted.", slug, slug)
			overallValid = false
		}
	} else {
		jailActiveMsg = "Fail2Ban is not running or active jails could not be verified."
		overallValid = false
	}
	checks = append(checks, models.SetupValidationCheck{Name: "Fail2Ban Jails Active", Passed: jailsActive, Message: jailActiveMsg})

	return &models.SetupValidationResponse{
		Checks:       checks,
		OverallValid: overallValid,
	}, nil
}

// ValidateRemoval validates Nginx cleanup.
func (s *SetupService) ValidateRemoval(domainName string) (*models.RemovalValidationResponse, error) {
	if s.demoMode {
		return &models.RemovalValidationResponse{
			Checks: []models.SetupValidationCheck{
				{Name: "Nginx Config Cleaned", Passed: true, Message: "Nginx config files no longer reference protection variables for " + domainName},
			},
			OverallValid: true,
		}, nil
	}

	slug := domainToSlug(domainName)
	varName := domainToVarName(slug)
	overallValid := true
	checks := []models.SetupValidationCheck{}

	dirsToScan := []string{}
	hostDirs := []string{
		"/host/nginx-sites-enabled",
		"/host/nginx-conf-d",
		"/host/nginx-sites-available",
		"/host/nginx-config",
	}
	stdDirs := []string{
		"/etc/nginx/sites-enabled",
		"/etc/nginx/conf.d",
		"/etc/nginx/sites-available",
		"/etc/nginx",
	}

	for i, hostDir := range hostDirs {
		if _, err := os.Stat(hostDir); err == nil {
			dirsToScan = append(dirsToScan, hostDir)
		} else {
			if _, err := os.Stat(stdDirs[i]); err == nil {
				dirsToScan = append(dirsToScan, stdDirs[i])
			}
		}
	}

	if len(dirsToScan) == 0 {
		dirsToScan = stdDirs
	}

	referencingFiles := []string{}

	for _, dir := range dirsToScan {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, f := range files {
			if f.IsDir() {
				continue
			}
			filePath := filepath.Join(dir, f.Name())
			if strings.Contains(f.Name(), "_blocked.conf") {
				continue
			}
			contentBytes, err := os.ReadFile(filePath)
			if err != nil {
				continue
			}
			content := string(contentBytes)

			if strings.Contains(content, slug+"_rate_limit") ||
				strings.Contains(content, varName+"_blocked") ||
				strings.Contains(content, slug+"_blocked.conf") {

				displayFile := filePath
				if strings.HasPrefix(filePath, "/host/") {
					displayFile = "/etc/" + strings.TrimPrefix(filePath, "/host/nginx-")
				}
				referencingFiles = append(referencingFiles, displayFile)
			}
		}
	}

	passed := len(referencingFiles) == 0
	msg := "No references to domain protection rules found in active Nginx configs."
	if !passed {
		msg = fmt.Sprintf("Nginx configuration files still reference protection: %s. Please remove the snippet from these files.", strings.Join(referencingFiles, ", "))
		overallValid = false
	}

	checks = append(checks, models.SetupValidationCheck{
		Name:    "Nginx Config Cleaned",
		Passed:  passed,
		Message: msg,
	})

	return &models.RemovalValidationResponse{
		Checks:       checks,
		OverallValid: overallValid,
	}, nil
}
