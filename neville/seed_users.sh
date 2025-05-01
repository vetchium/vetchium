#!/bin/bash

set -e # Exit immediately if a command exits with a non-zero status.

echo "--- [Script] Fetching Postgres URI and Seeding Users ---"

# Fetch Postgres URI from user-specific namespace
EFFECTIVE_POSTGRES_URI=$(kubectl -n vetchium-devtest-$USER get secret postgres-app -o jsonpath='{.data.uri}' | base64 -d | sed 's/postgres-rw.vetchium-devtest-'$USER'/localhost/g')

if [ -z "$EFFECTIVE_POSTGRES_URI" ]; then
    echo "Error: [Script] Failed to retrieve Postgres URI from Kubernetes secret."
    exit 1
fi
echo "Using Postgres URI: $EFFECTIVE_POSTGRES_URI"

# Variables
NUM_USERS=${NUM_USERS:-100} # Read from environment or default to 100
HASHED_PW='$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK'
# Updated INSERT to include ALL missing non-null columns with defaults
INSERT_SQL_TPL="INSERT INTO hub_users (id, full_name, handle, email, password_hash, state, tier, resident_country_code, resident_city, preferred_language, short_bio, long_bio, created_at) VALUES (gen_random_uuid(), '%s', '%s', '%s', '%s', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'USA', 'Test City', 'en', 'Default short bio for test user.', 'Default long bio for test user.', NOW()) ON CONFLICT (handle) DO NOTHING;"

echo "Seeding $NUM_USERS users..."

for i in $(seq 1 $NUM_USERS); do
    handle="hubuser$i"
    email="hubuser$i@example.com"
    full_name="Hub User $i" # Note: No extra quoting needed here for the variable itself
    
    # Use printf to handle insertion, including proper quoting for the SQL string itself
    sql=$(printf "$INSERT_SQL_TPL" "$full_name" "$handle" "$email" "$HASHED_PW")
    
    # Execute using psql -c
    psql "$EFFECTIVE_POSTGRES_URI" -qt -c "$sql" || { 
        echo "Error: [Script] psql command failed for user $i."
        echo "Attempted SQL: $sql"
        exit 1; 
    }
done

echo "--- [Script] User seeding attempt complete ---"
exit 0
