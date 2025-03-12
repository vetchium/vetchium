BEGIN;

-- Create test hub users
INSERT INTO hub_users (
    id, full_name, handle, email, password_hash,
    state, tier, resident_country_code, resident_city, preferred_language,
    short_bio, long_bio, created_at, updated_at
) VALUES
(
    '12345678-0019-0019-0019-000000000001',
    'Profile Test User 1',
    'profilepage_user1',
    'user1@profilepage-hub.example',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    'ACTIVE_HUB_USER',
    'FREE_HUB_USER',
    'IND',
    'Mumbai',
    'en',
    'Profile Test User 1 is experienced',
    'Profile Test User 1 was born in India and finished education at IIT Mumbai.',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
(
    '12345678-0019-0019-0019-000000000002',
    'Profile Test User 2',
    'profilepage_user2',
    'user2@profilepage-hub.example',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    'ACTIVE_HUB_USER',
    'FREE_HUB_USER',
    'IND',
    'Delhi',
    'en',
    'Profile Test User 2 is skilled',
    'Profile Test User 2 was born in India and finished education at Delhi University.',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
(
    '12345678-0019-0019-0019-000000000003',
    'Profile Test User 3',
    'profilepage_user3',
    'user3@profilepage-hub.example',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    'ACTIVE_HUB_USER',
    'FREE_HUB_USER',
    'IND',
    'Chennai',
    'en',
    'Profile Test User 3',
    'Profile Test User 3 is dedicated for profile picture end-to-end testing.',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
);

COMMIT;
