BEGIN;

-- Create test hub users
INSERT INTO hub_users (
    id, full_name, handle, email, password_hash,
    state, resident_country_code, preferred_language
) VALUES
(
    '12345678-0017-0017-0017-000000000001',
    'Work History Test User 1',
    'workhistory-user1',
    'user1@workhistory-hub.example',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    'ACTIVE_HUB_USER',
    'IND',
    'en'
),
(
    '12345678-0017-0017-0017-000000000002',
    'Work History Test User 2',
    'workhistory-user2',
    'user2@workhistory-hub.example',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    'ACTIVE_HUB_USER',
    'IND',
    'en'
);

-- Create test employers
INSERT INTO employers (
    id, client_id_type, employer_state, company_name,
    onboard_admin_email
) VALUES
(
    '12345678-0017-0017-0017-000000000003',
    'DOMAIN',
    'ONBOARDED',
    'WorkHistory Employer 1',
    'admin@workhistory-employer1.example'
),
(
    '12345678-0017-0017-0017-000000000004',
    'DOMAIN',
    'ONBOARDED',
    'WorkHistory Employer 2',
    'admin@workhistory-employer2.example'
);

-- Create domains for employers
INSERT INTO domains (
    id, domain_name, domain_state, employer_id
) VALUES
(
    '12345678-0017-0017-0017-000000000005',
    'workhistory-employer1.example',
    'VERIFIED',
    '12345678-0017-0017-0017-000000000003'
),
(
    '12345678-0017-0017-0017-000000000006',
    'workhistory-employer2.example',
    'VERIFIED',
    '12345678-0017-0017-0017-000000000004'
),
(
    '12345678-0017-0017-0017-000000000009',
    'non-onboarded-employer.example',
    'UNVERIFIED',
    NULL
);

-- Set primary domains for employers
INSERT INTO employer_primary_domains (
    employer_id, domain_id
) VALUES
(
    '12345678-0017-0017-0017-000000000003',
    '12345678-0017-0017-0017-000000000005'
),
(
    '12345678-0017-0017-0017-000000000004',
    '12345678-0017-0017-0017-000000000006'
);

-- Create some initial work history entries
INSERT INTO work_history (
    id, hub_user_id, employer_id, title, start_date, end_date, description
) VALUES
(
    '12345678-0017-0017-0017-000000000007',
    '12345678-0017-0017-0017-000000000001',
    '12345678-0017-0017-0017-000000000003',
    'Software Engineer',
    '2020-01-01',
    '2021-12-31',
    'Worked on various projects'
),
(
    '12345678-0017-0017-0017-000000000008',
    '12345678-0017-0017-0017-000000000001',
    '12345678-0017-0017-0017-000000000004',
    'Senior Engineer',
    '2022-01-01',
    NULL,
    'Currently working on cloud infrastructure'
);

COMMIT;
