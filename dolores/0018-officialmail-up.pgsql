BEGIN;

-- Create test hub users (one for each API test)
INSERT INTO hub_users (
    id, full_name, handle, email, password_hash,
    state, resident_country_code, resident_city, preferred_language,
    short_bio, long_bio, created_at, updated_at
) VALUES
-- User for add official email tests
(
    '12345678-0018-0018-0018-000000000001',
    'Add Email Test User',
    'addemailuser',
    'addemailuser@hub.example',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    'ACTIVE_HUB_USER',
    'IND',
    'Bangalore',
    'en',
    'Add Email Test User is organized',
    'Add Email Test User was born in India and finished education at IIIT Hyderabad and has 4 years as experience.',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
-- User for delete official email tests
(
    '12345678-0018-0018-0018-000000000002',
    'Delete Email Test User',
    'deleteemailuser',
    'deleteemailuser@hub.example',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    'ACTIVE_HUB_USER',
    'IND',
    'Chennai',
    'en',
    'Delete Email Test User is thorough',
    'Delete Email Test User was born in India and finished education at VIT Chennai and has 5 years as experience.',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
-- User for trigger verification tests
(
    '12345678-0018-0018-0018-000000000003',
    'Trigger Verification Test User',
    'triggeruser',
    'triggeruser@hub.example',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    'ACTIVE_HUB_USER',
    'IND',
    'Pune',
    'en',
    'Trigger Verification Test User is meticulous',
    'Trigger Verification Test User was born in India and finished education at Pune University and has 3 years as experience.',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
-- User for verify email tests
(
    '12345678-0018-0018-0018-000000000004',
    'Verify Email Test User',
    'verifyuser',
    'verifyuser@hub.example',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    'ACTIVE_HUB_USER',
    'IND',
    'Hyderabad',
    'en',
    'Verify Email Test User is precise',
    'Verify Email Test User was born in India and finished education at BITS Hyderabad and has 6 years as experience.',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
),
-- User for list emails tests
(
    '12345678-0018-0018-0018-000000000005',
    'List Emails Test User',
    'listemailsuser',
    'listemailsuser@hub.example',
    '$2a$10$p7Z/hRlt3ZZiz1IbPSJUiOualKbokFExYiWWazpQvfv660LqskAUK',
    'ACTIVE_HUB_USER',
    'IND',
    'Kolkata',
    'en',
    'List Emails Test User is systematic',
    'List Emails Test User was born in India and finished education at IIT Kharagpur and has 7 years as experience.',
    timezone('UTC'::text, now()),
    timezone('UTC'::text, now())
);

-- Create test employer
INSERT INTO employers (
    id, client_id_type, employer_state, company_name,
    onboard_admin_email
) VALUES
(
    '12345678-0018-0018-0018-000000000006',
    'DOMAIN',
    'ONBOARDED',
    'Official Mail Test Employer',
    'admin@officialmail.example'
);

-- Create domains for testing
INSERT INTO domains (
    id, domain_name, domain_state, employer_id
) VALUES
-- Verified domain
(
    '12345678-0018-0018-0018-000000000007',
    'officialmail.example',
    'VERIFIED',
    '12345678-0018-0018-0018-000000000006'
),
-- Unverified domain
(
    '12345678-0018-0018-0018-000000000008',
    'unverified.example',
    'UNVERIFIED',
    '12345678-0018-0018-0018-000000000006'
);

-- Set primary domain for employer
INSERT INTO employer_primary_domains (
    employer_id, domain_id
) VALUES
(
    '12345678-0018-0018-0018-000000000006',
    '12345678-0018-0018-0018-000000000007'
);

-- Create test data for add email tests (user1)
-- Empty initially, will test adding emails

-- Create test data for delete email tests (user2)
INSERT INTO hub_users_official_emails (
    hub_user_id,
    domain_id,
    official_email,
    verification_code,
    verification_code_expires_at,
    last_verified_at
) VALUES
(
    '12345678-0018-0018-0018-000000000002',
    '12345678-0018-0018-0018-000000000007',
    'delete.verified@officialmail.example',
    NULL,
    NULL,
    '2024-01-01 00:00:00'
),
(
    '12345678-0018-0018-0018-000000000002',
    '12345678-0018-0018-0018-000000000007',
    'delete.pending@officialmail.example',
    'DELETE123',
    timezone('UTC', now()) + interval '24 hours',
    NULL
);

-- Create test data for trigger verification tests (user3)
INSERT INTO hub_users_official_emails (
    hub_user_id,
    domain_id,
    official_email,
    verification_code,
    verification_code_expires_at,
    last_verified_at
) VALUES
(
    '12345678-0018-0018-0018-000000000003',
    '12345678-0018-0018-0018-000000000007',
    'trigger.recent@officialmail.example',
    NULL,
    NULL,
    timezone('UTC', now()) - interval '10 days'  -- Recent verification (10 days ago)
),
(
    '12345678-0018-0018-0018-000000000003',
    '12345678-0018-0018-0018-000000000007',
    'trigger.old@officialmail.example',
    NULL,
    NULL,
    timezone('UTC', now()) - interval '100 days'  -- Old verification (100 days ago)
);

-- Create test data for verify email tests (user4)
INSERT INTO hub_users_official_emails (
    hub_user_id,
    domain_id,
    official_email,
    verification_code,
    verification_code_expires_at,
    last_verified_at
) VALUES
(
    '12345678-0018-0018-0018-000000000004',
    '12345678-0018-0018-0018-000000000007',
    'verify.pending@officialmail.example',
    'VERIFY123',
    timezone('UTC', now()) + interval '24 hours',
    NULL
),
(
    '12345678-0018-0018-0018-000000000004',
    '12345678-0018-0018-0018-000000000007',
    'verify.expired@officialmail.example',
    'EXPIRED',
    timezone('UTC', now()) - interval '1 hour',
    NULL
);

-- Create test data for list emails tests (user5)
INSERT INTO hub_users_official_emails (
    hub_user_id,
    domain_id,
    official_email,
    verification_code,
    verification_code_expires_at,
    last_verified_at
) VALUES
(
    '12345678-0018-0018-0018-000000000005',
    '12345678-0018-0018-0018-000000000007',
    'list.verified@officialmail.example',
    NULL,
    NULL,
    '2024-01-01 00:00:00'
),
(
    '12345678-0018-0018-0018-000000000005',
    '12345678-0018-0018-0018-000000000007',
    'list.pending@officialmail.example',
    'LIST123',
    timezone('UTC', now()) + interval '24 hours',
    NULL
),
(
    '12345678-0018-0018-0018-000000000005',
    '12345678-0018-0018-0018-000000000007',
    'list.expired@officialmail.example',
    'EXPIRED',
    timezone('UTC', now()) - interval '1 hour',
    NULL
);

COMMIT;