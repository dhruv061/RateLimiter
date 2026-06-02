# Database Schema

SQLite is initialized automatically by the Go backend.

## Tables

- `users`: dashboard users, bcrypt password hashes, role, password-change status.
- `bans`: active and historical Fail2Ban records with IP, geo, jail, reason, timestamps, request counts, and active state.
- `whitelist`: trusted IPs with descriptions and audit-friendly ownership.
- `audit_logs`: user actions, targets, details, source IP, and timestamps.
- `settings`: key-value configuration for Nginx, Fail2Ban, dashboard behavior, and notifications.
- `traffic_stats`: minute/hour/day traffic rollups.

## Indexes

Indexes cover ban IP lookup, active bans, ban time, audit timestamps, traffic timestamps, and whitelist IP lookup.
