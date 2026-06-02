#!/bin/sh
set -eu

mkdir -p backups
stamp="$(date +%Y%m%d-%H%M%S)"
if [ -f data/dashboard.db ]; then
  cp data/dashboard.db "backups/dashboard-${stamp}.db"
  echo "Created backups/dashboard-${stamp}.db"
else
  echo "No database found at data/dashboard.db"
fi
