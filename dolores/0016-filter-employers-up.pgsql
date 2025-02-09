BEGIN;

-- Create test hub user
INSERT INTO hub_users (
    id, full_name, handle, email, password_hash,
    state, resident_country_code, preferred_language
) VALUES (
    '12345678-0016-0016-0016-000000000001',
    'Filter Employers Test User',
    'filter-employers',
    'user1@filter-employers.example',
    crypt('NewPassword123$', gen_salt('bf')),
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

COMMIT;
