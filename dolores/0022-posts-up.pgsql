BEGIN;

-- Create test hub users for posts functionality
INSERT INTO hub_users (
    id, full_name, handle, email, password_hash,
    state, tier, resident_country_code, resident_city, preferred_language,
    short_bio, long_bio, created_at
) VALUES
(
    '12345678-0022-0022-0022-000000000001',
    'Post Add Test User',
    'post-add-user',
    'add-user@0022-posts.example.com',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', -- NewPassword123$
    'ACTIVE_HUB_USER',
    'PAID_HUB_USER',
    'USA',
    'New York',
    'en',
    'Test user for adding posts',
    'This user is specifically for testing the add post functionality.',
    timezone('UTC'::text, now())
),
(
    '12345678-0022-0022-0022-000000000002',
    'Post Auth Test User',
    'post-auth-user',
    'auth-user@0022-posts.example.com',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', -- NewPassword123$
    'ACTIVE_HUB_USER',
    'PAID_HUB_USER',
    'GBR',
    'London',
    'en',
    'Test user for post authentication',
    'This user helps test authentication scenarios related to posts.',
    timezone('UTC'::text, now())
),
-- Users for GetUserPosts tests
(
    '12345678-0022-0022-0022-000000000003',
    'Get Posts User One',
    'get-user1',
    'get-user1@0022-posts.example.com',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', -- NewPassword123$
    'ACTIVE_HUB_USER',
    'PAID_HUB_USER',
    'CAN',
    'Toronto',
    'en',
    'User with multiple posts for get tests',
    'This user will have several posts to test fetching and pagination.',
    timezone('UTC'::text, now())
),
(
    '12345678-0022-0022-0022-000000000004',
    'Get Posts User Two',
    'get-user2',
    'get-user2@0022-posts.example.com',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', -- NewPassword123$
    'ACTIVE_HUB_USER',
    'PAID_HUB_USER',
    'FRA',
    'Paris',
    'fr',
    'Another user for get tests',
    'This user might have fewer posts or be used for handle lookup tests.',
    timezone('UTC'::text, now())
),
(
    '12345678-0022-0022-0022-000000000005',
    'Get Details User',
    'get-details-user',
    'get-details@0022-posts.example.com',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', -- NewPassword123$
    'ACTIVE_HUB_USER',
    'PAID_HUB_USER',
    'DEU',
    'Berlin',
    'de',
    'User for GetPostDetails tests',
    'Dedicated user for testing the single post details endpoint.',
    timezone('UTC'::text, now())
);

-- Add tags needed for GetUserPosts tests, ensuring they exist
-- Note: These tags should already exist from vetchium-tags.json
-- but we add them here in case they don't exist in test environment
INSERT INTO tags (id, name) VALUES
('productivity', 'Productivity'),
('innovation', 'Innovation'),
('golang', 'Go Programming Language'),
('technology', 'Technology')
ON CONFLICT (id) DO NOTHING;

COMMIT;
