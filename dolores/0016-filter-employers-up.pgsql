BEGIN;

-- Create test hub user
INSERT INTO hub_users (
    id, full_name, handle, email, password_hash,
    state, tier, resident_country_code, resident_city, preferred_language,
    short_bio, long_bio, created_at, updated_at
) VALUES (
    '12345678-0016-0016-0016-000000000001',
    'Filter Employers Test User',
    'filter-employers',
    'user1@filter-employers.example',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    'ACTIVE_HUB_USER',
    'FREE_HUB_USER',
    'IND',
    'Bangalore',
    'en',
    'Filter Employers Test User is detail-oriented',
    'Filter Employers Test User was born in India and finished education at IIIT Bangalore and has 3 years as experience.',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
);

-- Create test employers
INSERT INTO employers (
    id, client_id_type, employer_state, company_name,
    onboard_admin_email
) VALUES
(
    '12345678-0016-0016-0016-000000000002',
    'DOMAIN',
    'ONBOARDED',
    'Acme Corp',
    'admin@acme.example'
),
(
    '12345678-0016-0016-0016-000000000003',
    'DOMAIN',
    'ONBOARDED',
    'Beta Systems',
    'admin@beta.example'
);

-- Create domains for employers and domains without employers
INSERT INTO domains (
    id, domain_name, domain_state, employer_id
) VALUES
(
    '12345678-0016-0016-0016-000000000004',
    'acme.example',
    'VERIFIED',
    '12345678-0016-0016-0016-000000000002'
),
(
    '12345678-0016-0016-0016-000000000005',
    'beta.example',
    'VERIFIED',
    '12345678-0016-0016-0016-000000000003'
),
(
    '12345678-0016-0016-0016-000000000006',
    'domain-without-employer.example',
    'UNVERIFIED',
    NULL
);

-- Set primary domains for onboarded employers
INSERT INTO employer_primary_domains (
    employer_id, domain_id
) VALUES
(
    '12345678-0016-0016-0016-000000000002',
    '12345678-0016-0016-0016-000000000004'
),
(
    '12345678-0016-0016-0016-000000000003',
    '12345678-0016-0016-0016-000000000005'
);

COMMIT;
