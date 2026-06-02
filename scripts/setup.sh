#!/bin/sh
set -eu

cp -n .env.example .env 2>/dev/null || true
mkdir -p data logs backups
echo "Setup complete. Edit .env, then run: docker compose up -d --build"
