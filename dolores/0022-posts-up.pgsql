BEGIN;

-- Create test hub users for posts functionality
INSERT INTO hub_users (
    id, full_name, handle, email, password_hash,
    state, tier, resident_country_code, resident_city, preferred_language,
    short_bio, long_bio, created_at, updated_at
) VALUES
(
    '12345678-0022-0022-0022-000000000001',
    'Post Add Test User',
    'post-add-user',
    'add-user@0022-posts.example.com',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', -- NewPassword123$
    'ACTIVE_HUB_USER',
    'FREE_HUB_USER',
    'USA',
    'New York',
    'en',
    'Test user for adding posts',
    'This user is specifically for testing the add post functionality.',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
(
    '12345678-0022-0022-0022-000000000002',
    'Post Auth Test User',
    'post-auth-user',
    'auth-user@0022-posts.example.com',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', -- NewPassword123$
    'ACTIVE_HUB_USER',
    'FREE_HUB_USER',
    'GBR',
    'London',
    'en',
    'Test user for post authentication',
    'This user helps test authentication scenarios related to posts.',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
);

COMMIT;
