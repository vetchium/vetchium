BEGIN;

-- Clean up official emails
DELETE FROM hub_users_official_emails
WHERE hub_user_id IN (
    '12345678-0018-0018-0018-000000000001',
    '12345678-0018-0018-0018-000000000002',
    '12345678-0018-0018-0018-000000000003',
    '12345678-0018-0018-0018-000000000004',
    '12345678-0018-0018-0018-000000000005'
);

-- Clean up emails table
DELETE FROM emails
WHERE email_to && ARRAY[
    'add.new@officialmail.example',
    'delete.verified@officialmail.example',
    'delete.pending@officialmail.example',
    'trigger.recent@officialmail.example',
    'trigger.old@officialmail.example',
    'verify.pending@officialmail.example',
    'verify.expired@officialmail.example',
    'list.verified@officialmail.example',
    'list.pending@officialmail.example',
    'list.expired@officialmail.example'
];

-- Clean up primary domains
DELETE FROM employer_primary_domains
WHERE employer_id = '12345678-0018-0018-0018-000000000006';

-- Clean up domains
DELETE FROM domains
WHERE id IN (
    '12345678-0018-0018-0018-000000000007',
    '12345678-0018-0018-0018-000000000008'
);

-- Clean up employers
DELETE FROM employers
WHERE id = '12345678-0018-0018-0018-000000000006';

-- Clean up hub user tokens and related tables first
DELETE FROM hub_user_tfa_codes
WHERE tfa_token IN (
    SELECT token FROM hub_user_tokens
    WHERE hub_user_id IN (
        '12345678-0018-0018-0018-000000000001',
        '12345678-0018-0018-0018-000000000002',
        '12345678-0018-0018-0018-000000000003',
        '12345678-0018-0018-0018-000000000004',
        '12345678-0018-0018-0018-000000000005'
    )
);

DELETE FROM hub_user_tokens
WHERE hub_user_id IN (
    '12345678-0018-0018-0018-000000000001',
    '12345678-0018-0018-0018-000000000002',
    '12345678-0018-0018-0018-000000000003',
    '12345678-0018-0018-0018-000000000004',
    '12345678-0018-0018-0018-000000000005'
);

-- Finally clean up hub users
DELETE FROM hub_users
WHERE id IN (
    '12345678-0018-0018-0018-000000000001',
    '12345678-0018-0018-0018-000000000002',
    '12345678-0018-0018-0018-000000000003',
    '12345678-0018-0018-0018-000000000004',
    '12345678-0018-0018-0018-000000000005'
);

-- Cleanup for the dynamically created domain and employer from the new test case
-- These values must match what the test creates and the get_or_create_dummy_employer function produces.
DELETE FROM hub_users_official_emails WHERE official_email = 'contact@newlycreated.example.com';
DELETE FROM emails WHERE email_to && ARRAY['contact@newlycreated.example.com'];

-- Delete the primary domain entry first (if created - get_or_create_dummy_employer sets it as primary)
DELETE FROM employer_primary_domains
WHERE domain_id = (SELECT id FROM domains WHERE domain_name = 'newlycreated.example.com')
  AND employer_id = (SELECT id FROM employers WHERE company_name = 'newlycreated.example.com' AND employer_state = 'HUB_ADDED_EMPLOYER');

-- Then delete the domain
DELETE FROM domains WHERE domain_name = 'newlycreated.example.com';

-- Then delete the dummy employer
-- (get_or_create_dummy_employer uses domain name as company name for dummy record and sets state to HUB_ADDED_EMPLOYER)
DELETE FROM employers WHERE company_name = 'newlycreated.example.com' AND employer_state = 'HUB_ADDED_EMPLOYER';

COMMIT;
