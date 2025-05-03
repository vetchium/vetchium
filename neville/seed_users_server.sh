#!/bin/bash

set -e # Exit immediately if a command exits with a non-zero status.

echo "--- [Script] Fetching Postgres URI and Seeding Users via Server-Side Function ---"

# Check if VMUSER is set
if [ -z "$VMUSER" ]; then
    echo "Error: [Script] VMUSER environment variable is not set."
    exit 1
fi

# Check if VMADDR is set
if [ -z "$VMADDR" ]; then
    echo "Error: [Script] VMADDR environment variable is not set."
    exit 1
fi

# Fetch Postgres URI from user-specific namespace
EFFECTIVE_POSTGRES_URI=$(kubectl -n vetchium-devtest-$VMUSER get secret postgres-app -o jsonpath='{.data.uri}' | base64 -d | sed 's/postgres-rw.vetchium-devtest-'$VMUSER'/'$VMADDR'/g')

if [ -z "$EFFECTIVE_POSTGRES_URI" ]; then
    echo "Error: [Script] Failed to retrieve Postgres URI from Kubernetes secret."
    exit 1
fi
echo "Using Postgres URI: $EFFECTIVE_POSTGRES_URI"

# Variables
NUM_USERS=${NUM_USERS:-1000000} # Default to 1 million users
HASHED_PW='$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK'
START_INDEX=${START_INDEX:-1} # Start user index

echo "Seeding $NUM_USERS users starting from index $START_INDEX..."

# Create a temporary SQL file with the function
SQL_FILE=$(mktemp)
cat > "$SQL_FILE" << EOF
-- Create a function to generate users in bulk
CREATE OR REPLACE FUNCTION generate_test_users(start_idx INT, num_users INT, password_hash TEXT)
RETURNS VOID AS $$
DECLARE
    batch_size INT := 10000; -- Process in batches for better performance
    current_batch INT;
    remaining_users INT := num_users;
    current_start INT := start_idx;
BEGIN
    -- Notify start
    RAISE NOTICE 'Starting bulk user generation: % users from index %', num_users, start_idx;
    
    -- Process in batches
    WHILE remaining_users > 0 LOOP
        -- Calculate current batch size
        current_batch := LEAST(batch_size, remaining_users);
        
        RAISE NOTICE 'Processing batch of % users starting at index %', current_batch, current_start;
        
        -- Insert current batch using generate_series
        INSERT INTO hub_users (
            id, full_name, handle, email, password_hash, 
            state, tier, resident_country_code, resident_city, 
            preferred_language, short_bio, long_bio, created_at
        )
        SELECT 
            gen_random_uuid(),
            'Hub User ' || i,
            'hubuser' || i,
            'hubuser' || i || '@example.com',
            password_hash,
            'ACTIVE_HUB_USER',
            'FREE_HUB_USER',
            'USA',
            'Test City',
            'en',
            'Default short bio for test user.',
            'Default long bio for test user.',
            NOW()
        FROM generate_series(current_start, current_start + current_batch - 1) AS i
        ON CONFLICT (handle) DO NOTHING;
        
        -- Update counters
        remaining_users := remaining_users - current_batch;
        current_start := current_start + current_batch;
        
        -- Commit the current batch to avoid transaction size issues
        COMMIT;
        -- Start a new transaction
        BEGIN;
    END LOOP;
    
    RAISE NOTICE 'Bulk user generation completed successfully';
END;
$$ LANGUAGE plpgsql;

-- Execute the function
SELECT generate_test_users($START_INDEX, $NUM_USERS, '$HASHED_PW');

-- Clean up the function (optional)
-- DROP FUNCTION generate_test_users;
EOF

# Execute the SQL file
echo "Executing server-side user generation SQL function..."
psql "$EFFECTIVE_POSTGRES_URI" -f "$SQL_FILE"

# Remove the temporary SQL file
rm "$SQL_FILE"

echo "--- [Script] User seeding completed: $NUM_USERS users ---"
exit 0
