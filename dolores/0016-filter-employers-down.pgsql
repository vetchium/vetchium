BEGIN;

-- Clean up primary domains first
DELETE FROM employer_primary_domains
WHERE employer_id IN (
    '12345678-0016-0016-0016-000000000002',
    '12345678-0016-0016-0016-000000000003'
);

-- Clean up domains
DELETE FROM domains
WHERE id IN (
    '12345678-0016-0016-0016-000000000004',
    '12345678-0016-0016-0016-000000000005',
    '12345678-0016-0016-0016-000000000006'
);

-- Clean up employers
DELETE FROM employers
WHERE id IN (
    '12345678-0016-0016-0016-000000000002',
    '12345678-0016-0016-0016-000000000003'
);

-- Clean up hub user tokens and related tables first
DELETE FROM hub_user_tfa_codes
WHERE tfa_token IN (
    SELECT token FROM hub_user_tokens
    WHERE hub_user_id = '12345678-0016-0016-0016-000000000001'
);

DELETE FROM hub_user_tokens
WHERE hub_user_id = '12345678-0016-0016-0016-000000000001';

-- Finally clean up hub user
DELETE FROM hub_users
WHERE id = '12345678-0016-0016-0016-000000000001';

COMMIT;
