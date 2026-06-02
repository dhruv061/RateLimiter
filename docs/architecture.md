# Architecture

The dashboard is a single-VPS application that keeps the runtime small and operationally simple.

```text
Browser
  -> React/Vite dashboard
  -> Go/Gin API and WebSocket server
  -> SQLite database
  -> Host-mounted Nginx and Fail2Ban logs
```

The backend owns authentication, data persistence, audit logging, analytics aggregation, whitelist management, operational actions, and WebSocket fanout. The frontend is a dark-first SaaS dashboard shell with dense tables, right-side detail drawers, charts, and responsive navigation.

Persistent data is stored in host-mounted folders:

- `./data:/app/data`
- `./logs:/app/logs`
- `./backups:/app/backups`

No Redis, PostgreSQL, Kubernetes, Elasticsearch, or external queue is required.
