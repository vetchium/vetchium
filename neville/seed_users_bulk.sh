#!/bin/bash

set -e # Exit immediately if a command exits with a non-zero status.

echo "--- [Script] Bulk Creating Users for Load Testing ---"

# Check if PG_URI is set
if [ -z "$PG_URI" ]; then
    echo "Error: [Script] PG_URI environment variable is not set."
    exit 1
fi

# Check if NUM_USERS is set
if [ -z "$NUM_USERS" ]; then
    echo "Error: [Script] NUM_USERS environment variable is not set."
    exit 1
fi

echo "Using PostgreSQL URI: $PG_URI"
echo "Creating $NUM_USERS users..."

# First, load the bulk_create_users function
psql "$PG_URI" -c "\i /app/bulk_create_users.sql" || {
    echo "Error: [Script] Failed to load bulk_create_users function."
    exit 1
}

# Execute the bulk_create_users function
psql "$PG_URI" -c "SELECT bulk_create_users($NUM_USERS);" || {
    echo "Error: [Script] Failed to execute bulk_create_users function."
    exit 1
}

echo "--- [Script] User creation complete ---"
exit 0
