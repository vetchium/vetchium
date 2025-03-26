BEGIN;

-- Clean up education entries
DELETE FROM education
WHERE id IN (
    '12345678-0020-0020-0020-000000000006',
    '12345678-0020-0020-0020-000000000007'
) OR hub_user_id IN (
    '12345678-0020-0020-0020-000000000001',
    '12345678-0020-0020-0020-000000000002'
);

-- Clean up hub user tokens and related tables first
DELETE FROM hub_user_tfa_codes
WHERE tfa_token IN (
    SELECT token FROM hub_user_tokens
    WHERE hub_user_id IN (
        '12345678-0020-0020-0020-000000000001',
        '12345678-0020-0020-0020-000000000002'
    )
);

DELETE FROM hub_user_tokens
WHERE hub_user_id IN (
    '12345678-0020-0020-0020-000000000001',
    '12345678-0020-0020-0020-000000000002'
);

-- Clean up hub users
DELETE FROM hub_users
WHERE id IN (
    '12345678-0020-0020-0020-000000000001',
    '12345678-0020-0020-0020-000000000002'
);

-- Clean up institute domains
DELETE FROM institute_domains
WHERE domain IN (
    'annauniv.example',
    'psgtech.example',
    'stanford.example'
);

-- Clean up institutes
DELETE FROM institutes
WHERE id IN (
    '12345678-0020-0020-0020-000000000003',
    '12345678-0020-0020-0020-000000000004',
    '12345678-0020-0020-0020-000000000005'
);

COMMIT;
