#!/bin/sh
set -e

echo "Waiting for PostgreSQL to be ready..."

# Wait for PostgreSQL to be ready
until nc -z ${POSTGRESQL_CONTROL_HOST:-postgres} ${POSTGRESQL_CONTROL_PORT:-5432}; do
  echo "PostgreSQL is unavailable - sleeping"
  sleep 2
done

echo "PostgreSQL is up - executing migrations"

# Build database URL
DB_URL="postgresql://${POSTGRESQL_CONTROL_USER:-postgres}:${POSTGRESQL_CONTROL_PASSWORD:-password123}@${POSTGRESQL_CONTROL_HOST:-postgres}:${POSTGRESQL_CONTROL_PORT:-5432}/${POSTGRESQL_CONTROL_DB:-misviaticos_control}?sslmode=disable"

# Run migrations
migrate -path=/app/db/migrations -database "${DB_URL}" up

if [ $? -eq 0 ]; then
  echo "Migrations completed successfully"
else
  echo "Migration failed"
  exit 1
fi

# Run seeds
echo "Running database seeds..."
for seed_file in /app/db/seed/*.sql; do
  if [ -f "$seed_file" ]; then
    echo "Executing seed: $(basename $seed_file)"
    psql "${DB_URL}" -f "$seed_file"
    if [ $? -ne 0 ]; then
      echo "Warning: Seed $(basename $seed_file) failed, but continuing..."
    fi
  fi
done
echo "Seeds execution completed"

# Start the application
echo "Starting MisViaticos API..."
exec ./misviaticos-api
