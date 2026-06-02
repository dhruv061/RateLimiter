# PROJECT_GUIDE.md

# Fail2Ban & Nginx Rate Limit Dashboard

## Purpose

This repository contains a complete modern dashboard for managing:

- Nginx Rate Limiting
- Fail2Ban Bans
- Active Attack Monitoring
- Traffic Analytics
- IP Whitelisting
- Security Audit Logs

This file is the central entry point for developers and AI coding agents.

Before starting implementation, read the referenced documentation files.

---

# Documentation Structure

```text
docs/

├── PROJECT_GUIDE.md
├── design-principles.md
├── features.md
├── techstack.md
├── api-spec.md
├── database-schema.md
├── deployment.md
├── development-guide.md
└── architecture.md
```

---

# Documentation Overview

---

## design-principles.md

Purpose:

Defines how the application should look and behave.

Contains:

- UI Philosophy
- Light Theme Guidelines
- Dark Theme Guidelines
- Typography Rules
- Layout Rules
- Card Design
- Table Design
- Drawer Design
- Animation Rules
- User Experience Guidelines

Read this before building any UI.

All frontend implementation must follow this document.

---

## features.md

Purpose:

Defines what the application should do.

Contains:

- Dashboard Features
- Active Bans
- Ban History
- Live Requests
- Analytics
- Whitelist Management
- Audit Logs
- Settings
- Real-Time Updates

Read this before implementing any feature.

All feature development must align with this document.

---

## techstack.md

Purpose:

Defines approved technologies.

Contains:

- Frontend Stack
- Backend Stack
- Database Selection
- Docker Requirements
- Deployment Standards

No technology outside this document should be introduced without approval.

---

## architecture.md

Purpose:

System-level architecture.

Contains:

- Component Relationships
- Data Flow
- Authentication Flow
- WebSocket Flow
- Log Processing Flow

Use when building new modules.

---

## api-spec.md

Purpose:

Defines backend APIs.

Contains:

- REST Endpoints
- Request Models
- Response Models
- Authentication Requirements
- WebSocket Events

Frontend developers should reference this document.

---

## database-schema.md

Purpose:

Defines database structure.

Contains:

- Tables
- Relationships
- Indexes
- Migrations

Backend developers should reference this document.

---

## deployment.md

Purpose:

Production deployment instructions.

Contains:

- Docker Compose
- Environment Variables
- Reverse Proxy Setup
- SSL Setup
- Backup Strategy

Used during deployment.

---

## development-guide.md

Purpose:

Developer onboarding guide.

Contains:

- Local Setup
- Running Services
- Coding Standards
- Git Workflow
- Pull Request Rules

Read before contributing.

---

# Development Rules

---

## UI Development

Before creating any UI:

Read:

```text
docs/design-principles.md
```

Follow:

- Layout Standards
- Color Standards
- Component Standards
- Accessibility Standards

---

## Feature Development

Before building any feature:

Read:

```text
docs/features.md
```

Ensure:

- Feature exists in specification
- UI follows design principles
- API follows api-spec.md

---

## Backend Development

Before creating APIs:

Read:

```text
docs/api-spec.md
```

And:

```text
docs/database-schema.md
```

---

## Deployment Changes

Before modifying deployment:

Read:

```text
docs/deployment.md
```

---

# Repository Structure

```text
fail2ban-dashboard/

├── frontend/
├── backend/
├── docs/
├── docker/
├── scripts/
├── nginx/
├── data/
├── logs/
├── backups/
├── .env.example
├── docker-compose.yml
├── Makefile
└── README.md
```

---

# Frontend Structure

```text
frontend/

├── public/
├── src/

│   ├── app/
│   ├── pages/
│   ├── layouts/
│   ├── components/
│   │
│   ├── features/
│   │   ├── dashboard/
│   │   ├── bans/
│   │   ├── analytics/
│   │   ├── whitelist/
│   │   ├── audit/
│   │   └── settings/
│   │
│   ├── hooks/
│   ├── services/
│   ├── store/
│   ├── types/
│   ├── utils/
│   └── styles/

├── package.json
└── vite.config.ts
```

---

# Backend Structure

```text
backend/

├── cmd/
│   └── server/

├── internal/

│   ├── api/
│   ├── auth/
│   ├── bans/
│   ├── analytics/
│   ├── whitelist/
│   ├── audit/
│   ├── settings/
│   ├── websocket/
│   ├── logs/
│   ├── database/
│   └── middleware/

├── migrations/
├── configs/
├── pkg/

├── go.mod
└── main.go
```

---

# Docker Structure

```text
docker/

├── frontend/
│   └── Dockerfile

├── backend/
│   └── Dockerfile

└── nginx/
    └── Dockerfile
```

---

# Documentation Structure

```text
docs/

├── PROJECT_GUIDE.md
├── design-principles.md
├── features.md
├── techstack.md
├── architecture.md
├── api-spec.md
├── database-schema.md
├── deployment.md
└── development-guide.md
```

---

# Persistent Data Structure

Host Machine

```text
/opt/fail2ban-dashboard/

├── data/
├── logs/
├── backups/
├── nginx/
└── config/
```

---

# Docker Compose Structure

Containers:

```text
dashboard
```

Contains:

- React Frontend
- Go Backend

Optional:

```text
dashboard-nginx
```

For reverse proxy.

---

# Environment Variables

Store in:

```text
.env
```

Example:

```env
APP_PORT=8080

JWT_SECRET=CHANGE_ME

DATABASE_PATH=/app/data/dashboard.db

NGINX_LOG_PATH=/host/nginx/domain_access.log

FAIL2BAN_LOG_PATH=/host/fail2ban.log

BLOCK_FILE_PATH=/host/fail2ban_blocked.conf
```

---

# README Generation Requirements

Root README.md must include:

## Project Overview

Explain:

- What the dashboard does
- Key features
- Screenshots

---

## Requirements

- Docker
- Docker Compose

---

## Installation

Step-by-step setup.

---

## Configuration

Environment variables.

---

## Running

Commands:

```bash
docker compose up -d
```

---

## Updating

Commands:

```bash
git pull
docker compose build
docker compose up -d
```

---

## Backup

Explain:

- Database backup
- Config backup

---

## Troubleshooting

Common issues.

---

# Development Workflow

Step 1

Read:

```text
docs/techstack.md
```

Step 2

Read:

```text
docs/design-principles.md
```

Step 3

Read:

```text
docs/features.md
```

Step 4

Review:

```text
docs/api-spec.md
```

Step 5

Implement feature.

---

# Golden Rules

1. Simplicity over complexity.

2. Modern UI over traditional admin panels.

3. Dark mode first.

4. Mobile-friendly but desktop-focused.

5. Every feature must exist in features.md.

6. Every UI must follow design-principles.md.

7. Every technology choice must follow techstack.md.

8. Everything must run through Docker.

9. All persistent data must be mounted to the host machine.

10. Documentation is part of the product and must be maintained with the code.