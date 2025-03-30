BEGIN;

-- Create test hub users
INSERT INTO hub_users (
    id, full_name, handle, email, password_hash,
    state, tier, resident_country_code, resident_city, preferred_language,
    short_bio, long_bio, created_at, updated_at
) VALUES
(
    '12345678-0021-0021-0021-000000000001',
    'Achievement Add Test User',
    'add-achievement-user',
    'add-user@achievement-hub.example',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    'ACTIVE_HUB_USER',
    'FREE_HUB_USER',
    'IND',
    'Chennai',
    'en',
    'Test user for adding achievements',
    'This user is for testing the add achievement functionality.',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
(
    '12345678-0021-0021-0021-000000000002',
    'Achievement List Test User',
    'list-achievement-user',
    'list-user@achievement-hub.example',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    'ACTIVE_HUB_USER',
    'FREE_HUB_USER',
    'USA',
    'Boston',
    'en',
    'Test user for listing achievements',
    'This user is for testing the list achievement functionality.',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
(
    '12345678-0021-0021-0021-000000000003',
    'Achievement Delete Test User',
    'delete-achievement-user',
    'delete-user@achievement-hub.example',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    'ACTIVE_HUB_USER',
    'FREE_HUB_USER',
    'GBR',
    'London',
    'en',
    'Test user for deleting achievements',
    'This user is for testing the delete achievement functionality.',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
(
    '12345678-0021-0021-0021-000000000004',
    'Achievement Flow Test User',
    'flow-achievement-user',
    'flow-user@achievement-hub.example',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    'ACTIVE_HUB_USER',
    'FREE_HUB_USER',
    'CAN',
    'Toronto',
    'en',
    'Test user for achievement workflow',
    'This user is for testing the complete achievement workflow.',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
(
    '12345678-0021-0021-0021-000000000005',
    'Achievement Second User',
    'second-achievement-user',
    'second-user@achievement-hub.example',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    'ACTIVE_HUB_USER',
    'FREE_HUB_USER',
    'AUS',
    'Sydney',
    'en',
    'Secondary test user for achievements',
    'This user is for viewing other users achievements.',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
(
    '12345678-0021-0021-0021-000000000006',
    'Achievement Employer View User',
    'employer-view-achievement-user',
    'employer-view@achievement-hub.example',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    'ACTIVE_HUB_USER',
    'FREE_HUB_USER',
    'DEU',
    'Berlin',
    'en',
    'Test user for employer views',
    'This user is for testing employer viewing achievements.',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
);

-- Create some initial achievements
INSERT INTO achievements (
    id, hub_user_id, achievement_type, title, description, url, achieved_at, created_at, updated_at
) VALUES
(
    '12345678-0021-0021-0021-000000000010',
    '12345678-0021-0021-0021-000000000002',
    'PATENT',
    'Machine Learning Patent',
    'A patent for innovative ML algorithms',
    'https://patent.example.com/ml-innovation',
    '2022-06-15',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
(
    '12345678-0021-0021-0021-000000000011',
    '12345678-0021-0021-0021-000000000002',
    'PUBLICATION',
    'Research on AI Ethics',
    'Publication about ethical considerations in AI development',
    'https://journal.example.com/ai-ethics',
    '2023-03-10',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
(
    '12345678-0021-0021-0021-000000000012',
    '12345678-0021-0021-0021-000000000003',
    'CERTIFICATION',
    'AWS Solutions Architect',
    'Professional certification for AWS architecture',
    'https://aws.example.com/certification',
    '2021-11-20',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
(
    '12345678-0021-0021-0021-000000000013',
    '12345678-0021-0021-0021-000000000006',
    'PATENT',
    'Blockchain Security Patent',
    'Patent for innovative blockchain security mechanisms',
    'https://patent.example.com/blockchain-security',
    '2020-08-05',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
(
    '12345678-0021-0021-0021-000000000014',
    '12345678-0021-0021-0021-000000000006',
    'PUBLICATION',
    'Research on Quantum Computing',
    'Publication about quantum computing applications',
    'https://journal.example.com/quantum-computing',
    '2022-01-15',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
);

-- Create test employer for org user testing
INSERT INTO employers (
    id, client_id_type, employer_state, company_name, onboard_admin_email
) VALUES 
(
    '12345678-0021-0021-0021-000000000020',
    'DOMAIN',
    'ONBOARDED',
    'Achievement Test Employer',
    'admin@achievement-employer.example'
);

-- Create domain for test employer
INSERT INTO domains (
    id, domain_name, domain_state, employer_id
) VALUES
(
    '12345678-0021-0021-0021-000000000021',
    'achievement-employer.example',
    'VERIFIED',
    '12345678-0021-0021-0021-000000000020'
);

-- Set primary domain for employer
INSERT INTO employer_primary_domains (
    employer_id, domain_id
) VALUES
(
    '12345678-0021-0021-0021-000000000020',
    '12345678-0021-0021-0021-000000000021'
);

-- Create org users for testing
INSERT INTO org_users (
    id, email, name, password_hash, org_user_roles, org_user_state, employer_id
) VALUES
(
    '12345678-0021-0021-0021-000000000022',
    'admin@achievement-employer.example',
    'Achievement Employer Admin',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    ARRAY['ADMIN']::org_user_roles[],
    'ACTIVE_ORG_USER',
    '12345678-0021-0021-0021-000000000020'
);

COMMIT;
