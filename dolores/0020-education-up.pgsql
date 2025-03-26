BEGIN;

-- Create test hub users
INSERT INTO hub_users (
    id, full_name, handle, email, password_hash,
    state, tier, resident_country_code, resident_city, preferred_language,
    short_bio, long_bio, created_at, updated_at
) VALUES
(
    '12345678-0020-0020-0020-000000000001',
    'Education Test User 1',
    'education-user1',
    'user1@education-hub.example',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    'ACTIVE_HUB_USER',
    'FREE_HUB_USER',
    'IND',
    'Chennai',
    'en',
    'Education Test User 1 is a student',
    'Education Test User 1 was born in Tamil Nadu and is studying at Anna University.',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
(
    '12345678-0020-0020-0020-000000000002',
    'Education Test User 2',
    'education-user2',
    'user2@education-hub.example',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    'ACTIVE_HUB_USER',
    'FREE_HUB_USER',
    'IND',
    'Coimbatore',
    'en',
    'Education Test User 2 is a graduate',
    'Education Test User 2 was born in Tamil Nadu and graduated from PSG Tech College.',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
);

-- Create test institutes
INSERT INTO institutes (
    id, institute_name, logo_url, created_at, updated_at
) VALUES
(
    '12345678-0020-0020-0020-000000000003',
    'Anna University',
    NULL,
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
(
    '12345678-0020-0020-0020-000000000004',
    'PSG College of Technology',
    NULL,
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
(
    '12345678-0020-0020-0020-000000000005',
    'Stanford University',
    NULL,
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
);

-- Create institute domains
INSERT INTO institute_domains (
    domain, institute_id, created_at, updated_at
) VALUES
(
    'annauniv.example',
    '12345678-0020-0020-0020-000000000003',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
(
    'psgtech.example',
    '12345678-0020-0020-0020-000000000004',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
(
    'stanford.example',
    '12345678-0020-0020-0020-000000000005',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
);

-- Create some initial education entries
INSERT INTO education (
    id, hub_user_id, institute_id, degree, start_date, end_date, description
) VALUES
(
    '12345678-0020-0020-0020-000000000006',
    '12345678-0020-0020-0020-000000000001',
    '12345678-0020-0020-0020-000000000003',
    'Bachelor of Computer Science',
    '2018-01-01',
    '2022-12-31',
    'Specialized in Artificial Intelligence'
),
(
    '12345678-0020-0020-0020-000000000007',
    '12345678-0020-0020-0020-000000000002',
    '12345678-0020-0020-0020-000000000004',
    'Master of Computer Applications',
    '2019-01-01',
    NULL,
    'Currently pursuing with focus on Data Science'
);

COMMIT;
