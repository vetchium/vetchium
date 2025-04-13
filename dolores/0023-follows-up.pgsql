BEGIN;

-- Create test hub users for follow/unfollow functionality
INSERT INTO hub_users (
    id, full_name, handle, email, password_hash,
    state, tier, resident_country_code, resident_city, preferred_language,
    short_bio, long_bio, created_at, updated_at
) VALUES
-- User that will follow others
(
    '12345678-0023-0023-0023-000000000001',
    'Follow Test User One',
    'follow-user1',
    'follow-user1@0023-follow.example.com',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', -- NewPassword123$
    'ACTIVE_HUB_USER',
    'FREE_HUB_USER',
    'USA',
    'New York',
    'en',
    'Test user who follows others',
    'This user is used to test the follow functionality.',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
-- User to be followed
(
    '12345678-0023-0023-0023-000000000002',
    'Follow Test User Two',
    'follow-user2',
    'follow-user2@0023-follow.example.com',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', -- NewPassword123$
    'ACTIVE_HUB_USER',
    'FREE_HUB_USER',
    'GBR',
    'London',
    'en',
    'Test user to be followed',
    'This user will be followed by others in the tests.',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
-- User to test unauthorized access
(
    '12345678-0023-0023-0023-000000000003',
    'Unauthorized Follow User',
    'unauth-follow-user',
    'unauth-user@0023-follow.example.com',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', -- NewPassword123$
    'ACTIVE_HUB_USER',
    'FREE_HUB_USER',
    'CAN',
    'Toronto',
    'en',
    'Unauthorized test user',
    'This user tests unauthorized scenarios.',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
-- User that's already following another user at test start
(
    '12345678-0023-0023-0023-000000000004',
    'Preexisting Follow User',
    'preexisting-follow',
    'preexisting@0023-follow.example.com',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', -- NewPassword123$
    'ACTIVE_HUB_USER',
    'FREE_HUB_USER',
    'AUS',
    'Sydney',
    'en',
    'User who already follows another user',
    'This user has existing follow relationships for testing.',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
-- User to be followed by preexisting-follow
(
    '12345678-0023-0023-0023-000000000005',
    'Followee Test User',
    'followee-user',
    'followee@0023-follow.example.com',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', -- NewPassword123$
    'ACTIVE_HUB_USER',
    'FREE_HUB_USER',
    'DEU',
    'Berlin',
    'en',
    'User who is already followed',
    'This user is already followed by preexisting-follow user.',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
-- User with deleted account
(
    '12345678-0023-0023-0023-000000000006',
    'Deleted User Account',
    'deleted-user',
    'deleted@0023-follow.example.com',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', -- NewPassword123$
    'DELETED_HUB_USER', -- This user has a deleted account
    'FREE_HUB_USER',
    'JPN',
    'Tokyo',
    'en',
    'User with deleted account',
    'This user has a deleted account for testing follow limitations.',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
-- Non-existent handle user
(
    '12345678-0023-0023-0023-000000000007',
    'Non Existent Handle',
    'non-existent-handle',
    'non-existent@0023-follow.example.com',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK', -- NewPassword123$
    'ACTIVE_HUB_USER',
    'FREE_HUB_USER',
    'FRA',
    'Paris',
    'en',
    'User with handle to test non-existence',
    'This user is for testing handle that does not exist.',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
);

-- Create preexisting follow relationship
INSERT INTO following_relationships (consuming_hub_user_id, producing_hub_user_id) VALUES
('12345678-0023-0023-0023-000000000004', '12345678-0023-0023-0023-000000000005');

COMMIT;
