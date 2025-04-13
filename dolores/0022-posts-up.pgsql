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
    'FREE_HUB_USER',
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
    'FREE_HUB_USER',
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
    'FREE_HUB_USER',
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
    'FREE_HUB_USER',
    'FRA',
    'Paris',
    'fr',
    'Another user for get tests',
    'This user might have fewer posts or be used for handle lookup tests.',
    timezone('UTC'::text, now())
);

-- Add posts for get-user1 (note varying timestamps)
INSERT INTO posts (id, content, author_id, created_at, updated_at) VALUES
('post-g1-01', 'First post by get-user1', '12345678-0022-0022-0022-000000000003', timezone('UTC'::text, now()) - interval '5 days', timezone('UTC'::text, now()) - interval '5 days'),
('post-g1-02', 'Second post by get-user1, with tags', '12345678-0022-0022-0022-000000000003', timezone('UTC'::text, now()) - interval '4 days', timezone('UTC'::text, now()) - interval '4 days'),
('post-g1-03', 'Third post, updated recently', '12345678-0022-0022-0022-000000000003', timezone('UTC'::text, now()) - interval '3 days', timezone('UTC'::text, now()) - interval '1 day'),
('post-g1-04', 'Fourth post, newest', '12345678-0022-0022-0022-000000000003', timezone('UTC'::text, now()) - interval '2 days', timezone('UTC'::text, now()));

-- Add posts for get-user2
INSERT INTO posts (id, content, author_id, created_at, updated_at) VALUES
('post-g2-01', 'First post by get-user2', '12345678-0022-0022-0022-000000000004', timezone('UTC'::text, now()) - interval '1 hour', timezone('UTC'::text, now()) - interval '1 hour');

-- Add tags needed for GetUserPosts tests, ensuring they exist
INSERT INTO tags (name) VALUES
('pagination'),
('specific-test'),
('golang'),
('testing')
ON CONFLICT (name) DO NOTHING;

-- Map tags to posts using a JOIN for robustness
INSERT INTO post_tags (post_id, tag_id)
SELECT
    p.post_id,
    t.id
FROM (VALUES
    ('post-g1-02', 'golang'),
    ('post-g1-02', 'testing'),
    ('post-g1-03', 'pagination'),
    ('post-g2-01', 'specific-test')
) AS p(post_id, tag_name)
JOIN tags t ON p.tag_name = t.name;

COMMIT;
