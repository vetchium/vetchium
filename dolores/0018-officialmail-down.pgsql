BEGIN;

-- Clean up official emails
DELETE FROM hub_users_official_emails
WHERE hub_user_id IN (
    '12345678-0018-0018-0018-000000000001',
    '12345678-0018-0018-0018-000000000002'
);

-- Clean up emails table
DELETE FROM emails
WHERE email_to && ARRAY['user1@officialmail.example', 'user2@officialmail.example'];

-- Clean up primary domains
DELETE FROM employer_primary_domains
WHERE employer_id = '12345678-0018-0018-0018-000000000003';

-- Clean up domains
DELETE FROM domains
WHERE id = '12345678-0018-0018-0018-000000000004';

-- Clean up employers
DELETE FROM employers
WHERE id = '12345678-0018-0018-0018-000000000003';

-- Clean up hub user tokens and related tables first
DELETE FROM hub_user_tfa_codes
WHERE tfa_token IN (
    SELECT token FROM hub_user_tokens
    WHERE hub_user_id IN (
        '12345678-0018-0018-0018-000000000001',
        '12345678-0018-0018-0018-000000000002'
    )
);

DELETE FROM hub_user_tokens
WHERE hub_user_id IN (
    '12345678-0018-0018-0018-000000000001',
    '12345678-0018-0018-0018-000000000002'
);

-- Finally clean up hub users
DELETE FROM hub_users
WHERE id IN (
    '12345678-0018-0018-0018-000000000001',
    '12345678-0018-0018-0018-000000000002'
);

COMMIT;
