.PHONY: help dev build up down restart logs clean setup backend frontend

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ─── Development ──────────────────────────────────────────

setup: ## Initial project setup
	@echo "🚀 Setting up Fail2Ban Dashboard..."
	@cp -n .env.example .env 2>/dev/null || true
	@mkdir -p data logs backups
	@echo "✅ Setup complete. Edit .env and run 'make dev'"

dev: ## Run both frontend and backend in development mode
	@echo "Starting development servers..."
	@$(MAKE) -j2 backend frontend

backend: ## Run Go backend in development mode
	cd backend && go run cmd/server/main.go

frontend: ## Run React frontend in development mode
	cd frontend && npm run dev

# ─── Docker ───────────────────────────────────────────────

build: ## Build Docker images
	docker compose build

up: ## Start all services
	docker compose up -d

down: ## Stop all services
	docker compose down

restart: ## Restart all services
	docker compose restart

logs: ## View logs
	docker compose logs -f

# ─── Database ─────────────────────────────────────────────

backup: ## Backup database
	@./scripts/backup.sh

seed: ## Seed demo data
	@./scripts/seed-demo.sh

# ─── Cleanup ──────────────────────────────────────────────

clean: ## Remove build artifacts
	rm -rf frontend/dist
	rm -rf backend/tmp
	docker compose down --rmi local --volumes 2>/dev/null || true
