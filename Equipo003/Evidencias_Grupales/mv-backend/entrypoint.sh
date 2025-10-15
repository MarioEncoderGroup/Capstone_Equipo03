#!/bin/sh
set -e

echo "🚀 MisViaticos API - Starting..."
echo "================================"

# Database connection details
DB_USER="${POSTGRESQL_CONTROL_USER:-postgres}"
DB_PASSWORD="${POSTGRESQL_CONTROL_PASSWORD:-password123}"
DB_HOST="${POSTGRESQL_CONTROL_HOST:-localhost}"
DB_PORT="${POSTGRESQL_CONTROL_PORT:-5432}"
DB_NAME="${POSTGRESQL_CONTROL_DATABASE:-misviaticos_control}"

# Build connection string
DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

# Wait for PostgreSQL to be ready
echo "⏳ Waiting for PostgreSQL to be ready..."
until pg_isready -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" > /dev/null 2>&1; do
  echo "   PostgreSQL is unavailable - sleeping"
  sleep 2
done
echo "✅ PostgreSQL is ready!"

# Run migrations
echo ""
echo "📊 Running database migrations..."
if migrate -path /app/db/migrations -database "${DATABASE_URL}" up; then
  echo "✅ Migrations completed successfully"
else
  echo "❌ Migration failed!"
  exit 1
fi

# Run seeds (only if tables are empty)
echo ""
echo "🌱 Running database seeds..."
SEED_DIR="/app/db/seed"
if [ -d "$SEED_DIR" ]; then
  for seed_file in "$SEED_DIR"/*.sql; do
    if [ -f "$seed_file" ]; then
      echo "   Executing: $(basename $seed_file)"
      PGPASSWORD="${DB_PASSWORD}" psql -h "${DB_HOST}" -U "${DB_USER}" -d "${DB_NAME}" -f "$seed_file" -v ON_ERROR_STOP=1 || true
    fi
  done
  echo "✅ Seeds completed"
else
  echo "⚠️  No seed directory found, skipping..."
fi

# Start the API
echo ""
echo "🎯 Starting MisViaticos API..."
echo "================================"
exec /app/misviaticos-api
