# API Specification

All responses use:

```json
{
  "success": true,
  "data": {}
}
```

Authentication uses JWT bearer tokens.

## Public

- `GET /health`
- `POST /api/auth/login`

## Authenticated

- `GET /api/auth/me`
- `POST /api/auth/change-password`
- `POST /api/auth/logout`
- `GET /api/dashboard/stats`
- `GET /api/dashboard/system-status`
- `GET /api/dashboard/attack-status`
- `GET /api/bans/active`
- `GET /api/bans/history`
- `GET /api/bans/top-offenders`
- `GET /api/bans/:id`
- `POST /api/bans/:id/unban`
- `POST /api/bans/bulk-unban`
- `GET /api/whitelist`
- `POST /api/whitelist`
- `DELETE /api/whitelist/:id`
- `GET /api/whitelist/export`
- `GET /api/settings`
- `PUT /api/settings`
- `GET /api/settings/validate`
- `GET /api/audit-logs`
- `GET /api/audit-logs/export`
- `GET /api/analytics/traffic-trends`
- `GET /api/analytics/countries`
- `GET /api/analytics/top-offenders`
- `GET /api/live-requests`
- `POST /api/system/nginx/validate`
- `POST /api/system/nginx/reload`
- `POST /api/system/nginx/restart`
- `POST /api/system/fail2ban/reload`
- `POST /api/system/fail2ban/restart`
- `POST /api/system/fail2ban/sync-bans`
- `GET /api/reports/security`

## WebSocket

- `GET /api/ws`

Events are typed envelopes:

```json
{
  "type": "dashboard.stats",
  "payload": {},
  "time": "2026-06-02T00:00:00Z"
}
```
