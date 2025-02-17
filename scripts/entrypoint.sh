#!/bin/sh

# Запуск миграций
echo "Running migrations..."
goose -dir /app/migrations postgres "host=$POSTGRES_HOST user=$POSTGRES_USER password=$POSTGRES_PASSWORD dbname=$POSTGRES_DB port=$POSTGRES_PORT sslmode=disable" up

# Запуск приложения
echo "Starting application..."
exec ./main 