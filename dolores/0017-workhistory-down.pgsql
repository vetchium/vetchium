BEGIN;

-- Clean up work history entries
DELETE FROM work_history
WHERE id IN (
    '12345678-0017-0017-0017-000000000007',
    '12345678-0017-0017-0017-000000000008'
);

-- Clean up primary domains
DELETE FROM employer_primary_domains
WHERE employer_id IN (
    '12345678-0017-0017-0017-000000000003',
    '12345678-0017-0017-0017-000000000004'
);

-- Clean up domains
DELETE FROM domains
WHERE id IN (
    '12345678-0017-0017-0017-000000000005',
    '12345678-0017-0017-0017-000000000006'
);

-- Clean up employers
DELETE FROM employers
WHERE id IN (
    '12345678-0017-0017-0017-000000000003',
    '12345678-0017-0017-0017-000000000004'
);

-- Clean up hub user tokens and related tables first
DELETE FROM hub_user_tfa_codes
WHERE tfa_token IN (
    SELECT token FROM hub_user_tokens
    WHERE hub_user_id IN (
        '12345678-0017-0017-0017-000000000001',
        '12345678-0017-0017-0017-000000000002'
    )
);

DELETE FROM hub_user_tokens
WHERE hub_user_id IN (
    '12345678-0017-0017-0017-000000000001',
    '12345678-0017-0017-0017-000000000002'
);

-- Finally clean up hub users
DELETE FROM hub_users
WHERE id IN (
    '12345678-0017-0017-0017-000000000001',
    '12345678-0017-0017-0017-000000000002'
);

COMMIT;
