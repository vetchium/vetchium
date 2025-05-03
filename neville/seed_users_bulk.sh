#!/bin/bash

set -e # Exit immediately if a command exits with a non-zero status.

echo "--- [Script] Fetching Postgres URI and Seeding Users in Bulk ---"

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
BATCH_SIZE=${BATCH_SIZE:-10000} # Process users in batches
HASHED_PW='$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK'

echo "Seeding $NUM_USERS users in batches of $BATCH_SIZE..."

# Create temporary directory for SQL files
TMP_DIR=$(mktemp -d)
echo "Created temporary directory: $TMP_DIR"

# Function to generate SQL file for a batch of users
generate_batch_sql() {
    local start_idx=$1
    local end_idx=$2
    local batch_file="$TMP_DIR/batch_${start_idx}_${end_idx}.sql"
    
    echo "BEGIN;" > "$batch_file"
    
    for i in $(seq $start_idx $end_idx); do
        handle="hubuser$i"
        email="hubuser$i@example.com"
        full_name="Hub User $i"
        
        echo "INSERT INTO hub_users (id, full_name, handle, email, password_hash, state, tier, resident_country_code, resident_city, preferred_language, short_bio, long_bio, created_at) VALUES (gen_random_uuid(), '$full_name', '$handle', '$email', '$HASHED_PW', 'ACTIVE_HUB_USER', 'FREE_HUB_USER', 'USA', 'Test City', 'en', 'Default short bio for test user.', 'Default long bio for test user.', NOW()) ON CONFLICT (handle) DO NOTHING;" >> "$batch_file"
    done
    
    echo "COMMIT;" >> "$batch_file"
    echo "Generated batch file for users $start_idx-$end_idx"
}

# Generate and execute batches
total_batches=$(( (NUM_USERS + BATCH_SIZE - 1) / BATCH_SIZE ))
echo "Will process $total_batches batches"

for batch in $(seq 1 $total_batches); do
    start_idx=$(( (batch - 1) * BATCH_SIZE + 1 ))
    end_idx=$(( batch * BATCH_SIZE ))
    
    # Ensure we don't exceed NUM_USERS
    if [ $end_idx -gt $NUM_USERS ]; then
        end_idx=$NUM_USERS
    fi
    
    echo "Processing batch $batch/$total_batches (users $start_idx-$end_idx)..."
    generate_batch_sql $start_idx $end_idx
    
    # Execute the batch file
    psql "$EFFECTIVE_POSTGRES_URI" -f "$TMP_DIR/batch_${start_idx}_${end_idx}.sql" || {
        echo "Error: [Script] psql command failed for batch $batch."
        exit 1
    }
    
    # Remove the batch file after processing
    rm "$TMP_DIR/batch_${start_idx}_${end_idx}.sql"
    
    echo "Batch $batch completed"
done

# Clean up temporary directory
rmdir "$TMP_DIR"
echo "Removed temporary directory: $TMP_DIR"

echo "--- [Script] User seeding completed: $NUM_USERS users ---"
exit 0
