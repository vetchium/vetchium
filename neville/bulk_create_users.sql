-- Function to bulk create users for load testing
-- This function creates users from user1@example.com to userN@example.com
CREATE OR REPLACE FUNCTION bulk_create_users(
    p_num_users INTEGER,
    p_password_hash TEXT DEFAULT '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK'
) RETURNS INTEGER AS $$
DECLARE
    v_batch_size INTEGER := 1000; -- Process in batches for better performance
    v_batches INTEGER;
    v_remaining INTEGER;
    v_start INTEGER;
    v_end INTEGER;
    v_inserted INTEGER := 0;
    v_batch INTEGER := 0;
BEGIN
    -- Calculate number of batches
    v_batches := p_num_users / v_batch_size;
    v_remaining := p_num_users % v_batch_size;
    
    -- Process full batches
    FOR v_batch IN 0..v_batches-1 LOOP
        v_start := v_batch * v_batch_size + 1;
        v_end := (v_batch + 1) * v_batch_size;
        
        RAISE NOTICE 'Processing batch % of % (users % to %)', v_batch+1, v_batches, v_start, v_end;
        
        -- Insert batch of users
        WITH numbers AS (
            SELECT generate_series(v_start, v_end) AS i
        )
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
            'User ' || i, 
            'user' || i, 
            'user' || i || '@example.com', 
            p_password_hash, 
            'ACTIVE_HUB_USER', 
            'FREE_HUB_USER', 
            'USA', 
            'Test City', 
            'en', 
            'Default short bio for test user ' || i, 
            'Default long bio for test user ' || i, 
            NOW()
        FROM numbers
        ON CONFLICT (handle) DO NOTHING;
        
        GET DIAGNOSTICS v_inserted = v_inserted + ROW_COUNT;
    END LOOP;
    
    -- Process remaining users (if any)
    IF v_remaining > 0 THEN
        v_start := v_batches * v_batch_size + 1;
        v_end := v_start + v_remaining - 1;
        
        RAISE NOTICE 'Processing final batch (users % to %)', v_start, v_end;
        
        -- Insert remaining users
        WITH numbers AS (
            SELECT generate_series(v_start, v_end) AS i
        )
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
            'User ' || i, 
            'user' || i, 
            'user' || i || '@example.com', 
            p_password_hash, 
            'ACTIVE_HUB_USER', 
            'FREE_HUB_USER', 
            'USA', 
            'Test City', 
            'en', 
            'Default short bio for test user ' || i, 
            'Default long bio for test user ' || i, 
            NOW()
        FROM numbers
        ON CONFLICT (handle) DO NOTHING;
        
        GET DIAGNOSTICS v_inserted = v_inserted + ROW_COUNT;
    END IF;
    
    RETURN v_inserted;
END;
$$ LANGUAGE plpgsql;
