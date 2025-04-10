#!/bin/bash

set -e  # Exit immediately if a command exits with a non-zero status
set -a  # Export all variables sourced from .env.test

# Load environment
source .env.test

# Stop and remove existing containers
echo "Stopping any existing test containers..."
docker-compose --env-file .env.test -f docker-compose.test.yml down

# Start test database
echo "Starting test database container..."
docker-compose --env-file .env.test -f docker-compose.test.yml up -d

# Wait for database to become healthy
echo "Waiting for test database to become healthy..."
until docker inspect --format='{{.State.Health.Status}}' tarot-postgres-test 2>/dev/null | grep -q healthy; do
    echo -n "."
    sleep 2
done
echo -e "\nTest database is healthy!"

# Run integration tests
echo "Running integration tests..."
go test -v ./tests/integration/...

# Cleanup
echo "Stopping and removing test containers..."
docker-compose --env-file .env.test -f docker-compose.test.yml down -v
