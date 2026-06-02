package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"fail2ban-dashboard/internal/api/handlers"
	"fail2ban-dashboard/internal/api/middleware"
	"fail2ban-dashboard/internal/auth"
	"fail2ban-dashboard/internal/config"
	"fail2ban-dashboard/internal/database"
	"fail2ban-dashboard/internal/services"
	"fail2ban-dashboard/internal/websocket"
	"fail2ban-dashboard/pkg/response"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := database.New(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	if err := db.SeedDefaultAdmin(cfg.DefaultAdminUser, cfg.DefaultAdminPass); err != nil {
		log.Fatalf("failed to seed default admin: %v", err)
	}
	if err := db.SeedDefaultSettings(); err != nil {
		log.Fatalf("failed to seed default settings: %v", err)
	}
	if cfg.DemoMode {
		if err := services.NewDemoService(db).SeedDemoData(); err != nil {
			log.Fatalf("failed to seed demo data: %v", err)
		}
	}

	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiry)
	hub := websocket.NewHub()
	go hub.Run()
	go broadcastDashboardEvents(hub, db)

	router := setupRouter(db, jwtManager, hub)

	addr := ":" + cfg.AppPort
	log.Printf("%s listening on %s", cfg.AppName, addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}

func setupRouter(db *database.DB, jwtManager *auth.JWTManager, hub *websocket.Hub) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), middleware.CORSMiddleware())

	auditSvc := services.NewAuditService(db)
	dashboardSvc := services.NewDashboardService(db)
	banSvc := services.NewBanService(db)
	authHandler := handlers.NewAuthHandler(db, jwtManager)
	dashboardHandler := handlers.NewDashboardHandler(dashboardSvc)
	bansHandler := handlers.NewBansHandler(banSvc, auditSvc)
	whitelistHandler := handlers.NewWhitelistHandler(services.NewWhitelistService(db), auditSvc)
	settingsHandler := handlers.NewSettingsHandler(services.NewSettingsService(db), auditSvc)
	auditHandler := handlers.NewAuditHandler(auditSvc)
	analyticsHandler := handlers.NewAnalyticsHandler(dashboardSvc, banSvc)
	liveHandler := handlers.NewLiveHandler()
	systemHandler := handlers.NewSystemHandler()
	reportsHandler := handlers.NewReportsHandler(dashboardSvc, banSvc, auditSvc)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := router.Group("/api")
	{
		api.POST("/auth/login", authHandler.Login)
		api.GET("/ws", hub.HandleWebSocket)

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(jwtManager))
		{
			protected.GET("/auth/me", authHandler.GetMe)
			protected.POST("/auth/change-password", authHandler.ChangePassword)
			protected.POST("/auth/logout", authHandler.Logout)

			protected.GET("/dashboard/stats", dashboardHandler.GetStats)
			protected.GET("/dashboard/system-status", dashboardHandler.GetSystemStatus)
			protected.GET("/dashboard/attack-status", dashboardHandler.GetAttackStatus)

			protected.GET("/bans/active", bansHandler.GetActiveBans)
			protected.GET("/bans/history", bansHandler.GetBanHistory)
			protected.GET("/bans/top-offenders", bansHandler.GetTopOffenders)
			protected.GET("/bans/:id", bansHandler.GetBanDetail)
			protected.POST("/bans/:id/unban", bansHandler.UnbanIP)
			protected.POST("/bans/bulk-unban", bansHandler.BulkUnban)

			protected.GET("/whitelist", whitelistHandler.GetAll)
			protected.POST("/whitelist", whitelistHandler.Add)
			protected.DELETE("/whitelist/:id", whitelistHandler.Remove)
			protected.GET("/whitelist/export", whitelistHandler.Export)

			protected.GET("/settings", settingsHandler.GetSettings)
			protected.PUT("/settings", settingsHandler.UpdateSettings)
			protected.GET("/settings/validate", settingsHandler.ValidateConfig)

			protected.GET("/audit-logs", auditHandler.GetLogs)
			protected.GET("/audit-logs/export", auditHandler.Export)

			protected.GET("/analytics/traffic-trends", analyticsHandler.GetTrafficTrends)
			protected.GET("/analytics/countries", analyticsHandler.GetCountryStats)
			protected.GET("/analytics/top-offenders", analyticsHandler.GetTopOffenders)

			protected.GET("/live-requests", liveHandler.GetRequests)

			protected.POST("/system/nginx/validate", systemHandler.ValidateNginx)
			protected.POST("/system/nginx/reload", systemHandler.ReloadNginx)
			protected.POST("/system/nginx/restart", systemHandler.RestartNginx)
			protected.POST("/system/fail2ban/reload", systemHandler.ReloadFail2Ban)
			protected.POST("/system/fail2ban/restart", systemHandler.RestartFail2Ban)
			protected.POST("/system/fail2ban/sync-bans", systemHandler.SyncBans)

			protected.GET("/reports/security", reportsHandler.SecurityReport)

		}
	}

	hasFrontend := registerStaticFrontend(router)
	if !hasFrontend {
		router.NoRoute(func(c *gin.Context) {
			response.NotFound(c, "Endpoint not found")
		})
	}

	return router
}

func broadcastDashboardEvents(hub *websocket.Hub, db *database.DB) {
	dashboardSvc := services.NewDashboardService(db)
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if stats, err := dashboardSvc.GetStats(); err == nil {
			hub.BroadcastEvent("dashboard.stats", stats)
		}
		if attack, err := dashboardSvc.GetAttackStatus(); err == nil {
			hub.BroadcastEvent("security.attack_status", attack)
		}
		hub.BroadcastEvent("system.clients", map[string]int{"count": hub.ClientCount()})
	}
}

func registerStaticFrontend(router *gin.Engine) bool {
	publicDir := os.Getenv("PUBLIC_DIR")
	if publicDir == "" {
		publicDir = "./public"
	}

	indexPath := filepath.Join(publicDir, "index.html")
	if _, err := os.Stat(indexPath); err != nil {
		return false
	}

	router.Static("/assets", filepath.Join(publicDir, "assets"))
	router.NoRoute(func(c *gin.Context) {
		if len(c.Request.URL.Path) >= 5 && c.Request.URL.Path[:5] == "/api/" {
			response.NotFound(c, "Endpoint not found")
			return
		}
		c.File(indexPath)
	})
	return true
}
