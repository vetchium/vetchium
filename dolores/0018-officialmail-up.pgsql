BEGIN;

-- Create test hub users
INSERT INTO hub_users (
    id, full_name, handle, email, password_hash,
    state, resident_country_code, preferred_language
) VALUES
(
    '12345678-0018-0018-0018-000000000001',
    'Official Mail Test User 1',
    'officialmailuser1',
    'officialmailuser1@hub.example',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    'ACTIVE_HUB_USER',
    'IND',
    'en'
),
(
    '12345678-0018-0018-0018-000000000002',
    'Official Mail Test User 2',
    'officialmailuser2',
    'officialmailuser2@hub.example',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    'ACTIVE_HUB_USER',
    'IND',
    'en'
);

-- Create test employer
INSERT INTO employers (
    id, client_id_type, employer_state, company_name,
    onboard_admin_email
) VALUES
(
    '12345678-0018-0018-0018-000000000003',
    'DOMAIN',
    'ONBOARDED',
    'Official Mail Test Employer',
    'admin@officialmail.example'
);

-- Create domain for employer
INSERT INTO domains (
    id, domain_name, domain_state, employer_id
) VALUES
(
    '12345678-0018-0018-0018-000000000004',
    'officialmail.example',
    'VERIFIED',
    '12345678-0018-0018-0018-000000000003'
);

-- Set primary domain for employer
INSERT INTO employer_primary_domains (
    employer_id, domain_id
) VALUES
(
    '12345678-0018-0018-0018-000000000003',
    '12345678-0018-0018-0018-000000000004'
);

-- Create some initial official emails for testing
INSERT INTO hub_users_official_emails (
    hub_user_id,
    domain_id,
    official_email,
    verification_code,
    verification_code_expires_at,
    last_verified_at
) VALUES
(
    '12345678-0018-0018-0018-000000000001',
    '12345678-0018-0018-0018-000000000004',
    'user1@officialmail.example',
    NULL,
    NULL,
    '2024-01-01 00:00:00'
);

COMMIT;