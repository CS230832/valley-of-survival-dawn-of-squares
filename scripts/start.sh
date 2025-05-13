#!/bin/bash
set -e

POSTGRES_PASSWORD=apppassword POSTGRES_USER=appuser POSTGRES_DB=appdb docker-entrypoint.sh postgres &
echo "Waiting for Postgres to be ready..."
until pg_isready -h 127.0.0.1 -p 5432; do
  sleep 1
done

echo "Postgres is ready. Ensuring database setup is complete..."
sleep 2

echo "Starting game server..."

/app/server