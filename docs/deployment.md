# Deployment

## Requirements

- Docker
- Docker Compose

## Setup

```bash
./scripts/setup.sh
docker compose up -d --build
```

## Volumes

The compose file mounts:

- `./data:/app/data`
- `./logs:/app/logs`
- `./backups:/app/backups`
- `/var/log/nginx:/host/nginx:ro`
- `/var/log/fail2ban.log:/host/fail2ban.log:ro`
- `/etc/nginx/fail2ban_blocked.conf:/host/fail2ban_blocked.conf`

## Backup

```bash
./scripts/backup.sh
```
