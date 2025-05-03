-- Create a function to generate test users in batches if it doesn't exist
CREATE OR REPLACE FUNCTION create_user_batch(start_idx INTEGER, end_idx INTEGER, password_hash TEXT)
RETURNS VOID AS $$
BEGIN
    INSERT INTO hub_users (
        id, 
        full_name, 
        handle, 
        email, 
        password_hash, 
        state, 
        tier, 
        resident_country_code, 
        resident_city, 
        preferred_language, 
        short_bio, 
        long_bio, 
        created_at
    )
    SELECT 
        gen_random_uuid(), 
        'User ' || n, 
        'user' || n, 
        'user' || n || '@example.com', 
        password_hash, 
        'ACTIVE_HUB_USER', 
        'FREE_HUB_USER', 
        'USA', 
        'Test City', 
        'en', 
        'Default short bio for test user ' || n, 
        'Default long bio for test user ' || n, 
        NOW()
    FROM generate_series(start_idx, end_idx) AS n
    ON CONFLICT (handle) DO NOTHING;
    
    RAISE NOTICE 'Created users from % to %', start_idx, end_idx;
END;
$$ LANGUAGE plpgsql;

-- Call the function with the specified parameters
DO $$
DECLARE
    total_users INTEGER := $TOTAL_USERS;
    batch_size INTEGER := 10000;
    num_batches INTEGER;
    remainder INTEGER;
    i INTEGER;
    start_idx INTEGER;
    end_idx INTEGER;
    password_hash TEXT := '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK';
BEGIN
    -- Calculate batches
    num_batches := total_users / batch_size;
    remainder := total_users % batch_size;
    
    RAISE NOTICE 'Creating % users in % batches of % (plus % remainder)', 
                 total_users, num_batches, batch_size, remainder;
    
    -- Process full batches
    FOR i IN 0..num_batches-1 LOOP
        start_idx := i * batch_size + 1;
        end_idx := (i + 1) * batch_size;
        
        RAISE NOTICE 'Processing batch % of % (users % to %)', 
                     i+1, num_batches, start_idx, end_idx;
        
        PERFORM create_user_batch(start_idx, end_idx, password_hash);
    END LOOP;
    
    -- Process remainder if any
    IF remainder > 0 THEN
        start_idx := num_batches * batch_size + 1;
        end_idx := start_idx + remainder - 1;
        
        RAISE NOTICE 'Processing final batch (users % to %)', 
                     start_idx, end_idx;
        
        PERFORM create_user_batch(start_idx, end_idx, password_hash);
    END IF;
    
    RAISE NOTICE 'User creation complete';
END;
$$;

-- Drop the function when we're done (optional)
-- Uncomment the line below if you want to drop the function after use
-- DROP FUNCTION IF EXISTS create_user_batch;
