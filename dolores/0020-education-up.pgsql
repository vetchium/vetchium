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
    'user1-education',
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
    'user2-education',
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
),
(
    '12345678-0020-0020-0020-000000000010',
    'Education Test User 3',
    'user3-education',
    'user3@education-hub.example',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    'ACTIVE_HUB_USER',
    'FREE_HUB_USER',
    'USA',
    'Boston',
    'en',
    'Education Test User 3 is a professional',
    'Education Test User 3 studied in USA and works as a software engineer.',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
(
    '12345678-0020-0020-0020-000000000011',
    'Education List Test User',
    'list-education-user',
    'list-user@education-hub.example',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    'ACTIVE_HUB_USER',
    'FREE_HUB_USER',
    'USA',
    'New York',
    'en',
    'User for testing list education features',
    'This user is dedicated to testing list education functionality.',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
(
    '12345678-0020-0020-0020-000000000012',
    'Education Delete Test User',
    'delete-education-user',
    'delete-user@education-hub.example',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    'ACTIVE_HUB_USER',
    'FREE_HUB_USER',
    'CAN',
    'Toronto',
    'en',
    'User for testing delete education features',
    'This user is dedicated to testing delete education functionality.',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
(
    '12345678-0020-0020-0020-000000000013',
    'Education Flow Test User',
    'flow-education-user',
    'flow-user@education-hub.example',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    'ACTIVE_HUB_USER',
    'FREE_HUB_USER',
    'GBR',
    'London',
    'en',
    'User for testing complete education flow',
    'This user is dedicated to testing the complete education workflow.',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
(
    '12345678-0020-0020-0020-000000000030',
    'Education Org View Test User',
    'org-view-education-user',
    'org-view-user@education-hub.example',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    'ACTIVE_HUB_USER',
    'FREE_HUB_USER',
    'USA',
    'San Francisco',
    'en',
    'User for testing org view of education',
    'This user is dedicated to testing org user view of education.',
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
),
(
    '12345678-0020-0020-0020-000000000014',
    'MIT',
    NULL,
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
(
    '12345678-0020-0020-0020-000000000015',
    'Caltech',
    NULL,
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
(
    '12345678-0020-0020-0020-000000000016',
    'Princeton University',
    NULL,
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
(
    '12345678-0020-0020-0020-000000000017',
    'Yale University',
    NULL,
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
(
    '12345678-0020-0020-0020-000000000018',
    'Columbia University',
    NULL,
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
(
    '12345678-0020-0020-0020-000000000019',
    'UC Berkeley',
    NULL,
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
(
    '12345678-0020-0020-0020-000000000020',
    'Oxford University',
    NULL,
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
(
    '12345678-0020-0020-0020-000000000031',
    'Harvard University',
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
),
(
    'mit.example',
    '12345678-0020-0020-0020-000000000014',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
(
    'caltech.example',
    '12345678-0020-0020-0020-000000000015',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
(
    'princeton.example',
    '12345678-0020-0020-0020-000000000016',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
(
    'yale.example',
    '12345678-0020-0020-0020-000000000017',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
(
    'columbia.example',
    '12345678-0020-0020-0020-000000000018',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
(
    'berkeley.example',
    '12345678-0020-0020-0020-000000000019',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
(
    'oxford.example',
    '12345678-0020-0020-0020-000000000020',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
(
    'harvard.example',
    '12345678-0020-0020-0020-000000000031',
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
),
(
    '12345678-0020-0020-0020-000000000021',
    '12345678-0020-0020-0020-000000000011',
    '12345678-0020-0020-0020-000000000014',
    'Bachelor of Engineering',
    '2015-09-01',
    '2019-05-31',
    'Electrical Engineering'
),
(
    '12345678-0020-0020-0020-000000000022',
    '12345678-0020-0020-0020-000000000011',
    '12345678-0020-0020-0020-000000000015',
    'Master of Science',
    '2019-09-01',
    '2021-05-31',
    'Computer Engineering'
),
(
    '12345678-0020-0020-0020-000000000023',
    '12345678-0020-0020-0020-000000000012',
    '12345678-0020-0020-0020-000000000019',
    'Master of Computer Science',
    '2022-01-01',
    '2024-01-01',
    'Specialized in AI'
),
(
    '12345678-0020-0020-0020-000000000032',
    '12345678-0020-0020-0020-000000000030',
    '12345678-0020-0020-0020-000000000031',
    'Bachelor of Arts',
    '2016-09-01',
    '2020-05-31',
    'Economics'
),
(
    '12345678-0020-0020-0020-000000000033',
    '12345678-0020-0020-0020-000000000030',
    '12345678-0020-0020-0020-000000000018',
    'Master of Business Administration',
    '2021-09-01',
    '2023-05-31',
    'Finance'
);

-- Create test employer for org user testing
INSERT INTO employers (
    id, client_id_type, employer_state, company_name, onboard_admin_email
) VALUES 
(
    '12345678-0020-0020-0020-000000000034',
    'DOMAIN',
    'ONBOARDED',
    'Education Test Employer',
    'admin@edu-employer.example'
);

-- Create domain for test employer
INSERT INTO domains (
    id, domain_name, domain_state, employer_id
) VALUES
(
    '12345678-0020-0020-0020-000000000035',
    'edu-employer.example',
    'VERIFIED',
    '12345678-0020-0020-0020-000000000034'
);

-- Set primary domain for employer
INSERT INTO employer_primary_domains (
    employer_id, domain_id
) VALUES
(
    '12345678-0020-0020-0020-000000000034',
    '12345678-0020-0020-0020-000000000035'
);

-- Create org users for testing
INSERT INTO org_users (
    id, email, name, password_hash, org_user_roles, org_user_state, employer_id
) VALUES
(
    '12345678-0020-0020-0020-000000000036',
    'admin@edu-employer.example',
    'Education Employer Admin',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    ARRAY['ADMIN']::org_user_roles[],
    'ACTIVE_ORG_USER',
    '12345678-0020-0020-0020-000000000034'
);

COMMIT;
