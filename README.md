# Fail2Ban & Nginx Rate Limit Dashboard

A modern dark-first dashboard for Nginx rate limiting, Fail2Ban bans, active attack monitoring, traffic analytics, IP whitelisting, settings, reports, and audit logs.

## Features

- JWT authentication with a seeded admin account.
- Active ban management, ban history, unban actions, and whitelist management.
- Real-time dashboard updates over WebSockets.
- Traffic trends, top offenders, country analytics, and live request monitoring.
- SQLite persistence with host-mounted data, logs, and backups.
- Docker Compose deployment for a single VPS.

## Requirements

- Docker
- Docker Compose

## Installation

```bash
./scripts/setup.sh
```

Edit `.env`, especially `JWT_SECRET` and default admin values.

## Running

```bash
docker compose up -d --build
```

Open:

```text
http://localhost:8080
```

Default demo credentials come from `.env`:

```text
admin / admin
```

Change the password after first login.

## Configuration

Key environment variables:

- `APP_PORT`: dashboard port.
- `JWT_SECRET`: signing secret for JWT tokens.
- `DATABASE_PATH`: SQLite database path inside the container.
- `NGINX_ACCESS_LOG`: mounted Nginx access log path.
- `NGINX_ERROR_LOG`: mounted Nginx error log path.
- `FAIL2BAN_LOG`: mounted Fail2Ban log path.
- `BLOCK_FILE_PATH`: mounted Fail2Ban/Nginx block include file.
- `DEMO_MODE`: seed demo data when true.

## Updating

```bash
git pull
docker compose build
docker compose up -d
```

## Backup

```bash
./scripts/backup.sh
```

Backups are written to `backups/`.

## Development

Backend:

```bash
cd backend
go run cmd/server/main.go
```

Frontend:

```bash
cd frontend
npm install
npm run dev
```

## Troubleshooting

- If login fails on a fresh database, check `DEFAULT_ADMIN_USER` and `DEFAULT_ADMIN_PASS`.
- If log-derived data is empty, verify host log mounts in `docker-compose.yml`.
- If the app cannot write data, check ownership and permissions for `data/`, `logs/`, and `backups/`.
- If WebSocket status shows polling, verify `/api/ws` is reachable through the reverse proxy.
