# Development Guide

## Backend

```bash
cd backend
go run cmd/server/main.go
```

## Frontend

```bash
cd frontend
npm install
npm run dev
```

The Vite dev server proxies `/api` and `/health` to `localhost:8080`.

## Docker

```bash
docker compose up -d --build
```
