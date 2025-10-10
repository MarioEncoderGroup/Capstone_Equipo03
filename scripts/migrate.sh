#!/bin/sh
# Script para ejecutar migraciones de MisViaticos
# Ejecuta automaticamente al levantar Docker Compose

set -e

echo "Esperando a que PostgreSQL este listo..."

# Esperar a que PostgreSQL este disponible
until PGPASSWORD=$POSTGRES_PASSWORD psql -h "$POSTGRES_HOST" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c '\q' 2>/dev/null; do
  echo "PostgreSQL no esta listo - esperando..."
  sleep 2
done

echo "PostgreSQL esta listo!"
echo ""

echo "Ejecutando migraciones de la base de datos de control..."
echo "Host: $POSTGRES_HOST"
echo "Database: $POSTGRES_DB"
echo "User: $POSTGRES_USER"
echo ""

# Construir DATABASE_URL
DATABASE_URL="postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable"

# Ejecutar migraciones usando golang-migrate
migrate -path=/migrations -database "${DATABASE_URL}" up

if [ $? -eq 0 ]; then
  echo "Migraciones ejecutadas exitosamente!"
else
  echo "Error ejecutando migraciones"
  exit 1
fi

echo ""
echo "Ejecutando scripts de inicializacion..."

# Ejecutar scripts de inicializacion (parametros, etc.)
for script in /init/*.sql; do
  if [ -f "$script" ]; then
    echo "Ejecutando: $(basename $script)"
    PGPASSWORD=$POSTGRES_PASSWORD psql -h "$POSTGRES_HOST" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -f "$script"
  fi
done

echo ""
echo "Inicializacion de base de datos completada!"
