BEGIN;

-- Create test hub users for comments functionality
INSERT INTO hub_users (
    id, full_name, handle, email, password_hash,
    state, tier, resident_country_code, resident_city, preferred_language,
    short_bio, long_bio, created_at
) VALUES
(
    '12345678-0028-0028-0028-000000000001',
    'Post Author User',
    'post-author-user',
    'post-author@0028-comments.example.com',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', -- NewPassword123$
    'ACTIVE_HUB_USER',
    'FREE_HUB_USER',
    'USA',
    'New York',
    'en',
    'Test user for post authoring',
    'This user creates posts for comment testing.',
    timezone('UTC'::text, now())
),
(
    '12345678-0028-0028-0028-000000000002',
    'Commenter User',
    'commenter-user',
    'commenter@0028-comments.example.com',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', -- NewPassword123$
    'ACTIVE_HUB_USER',
    'FREE_HUB_USER',
    'GBR',
    'London',
    'en',
    'Test user for commenting',
    'This user adds comments to posts.',
    timezone('UTC'::text, now())
),
(
    '12345678-0028-0028-0028-000000000003',
    'Other User',
    'other-user',
    'other-user@0028-comments.example.com',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', -- NewPassword123$
    'ACTIVE_HUB_USER',
    'FREE_HUB_USER',
    'CAN',
    'Toronto',
    'en',
    'Another test user',
    'This user is used for various test scenarios.',
    timezone('UTC'::text, now())
),
(
    '12345678-0028-0028-0028-000000000004',
    'Disable Comments User',
    'disable-comments-user',
    'disable-comments@0028-comments.example.com',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', -- NewPassword123$
    'ACTIVE_HUB_USER',
    'FREE_HUB_USER',
    'FRA',
    'Paris',
    'fr',
    'User for disable comments tests',
    'This user tests disabling comments functionality.',
    timezone('UTC'::text, now())
),
(
    '12345678-0028-0028-0028-000000000005',
    'Delete Comment User',
    'delete-comment-user',
    'delete-comment@0028-comments.example.com',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', -- NewPassword123$
    'ACTIVE_HUB_USER',
    'FREE_HUB_USER',
    'DEU',
    'Berlin',
    'de',
    'User for delete comment tests',
    'This user tests deleting comments as post author.',
    timezone('UTC'::text, now())
),
(
    '12345678-0028-0028-0028-000000000006',
    'Delete My Comment User',
    'delete-my-comment-user',
    'delete-my-comment@0028-comments.example.com',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', -- NewPassword123$
    'ACTIVE_HUB_USER',
    'FREE_HUB_USER',
    'AUS',
    'Sydney',
    'en',
    'User for delete my comment tests',
    'This user tests deleting their own comments.',
    timezone('UTC'::text, now())
);

-- Add tags needed for comment tests
INSERT INTO tags (name) VALUES
('0028-comments'),
('0028-test'),
('0028-disable'),
('0028-delete'),
('0028-get'),
('0028-my'),
('0028-enable')
ON CONFLICT (name) DO NOTHING;

COMMIT;